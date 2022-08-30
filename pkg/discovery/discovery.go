package discovery

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"sync"

	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/apic"
	"github.com/Axway/agent-sdk/pkg/util/log"
	"github.com/Axway/agents-layer7/pkg/client"
	"github.com/Axway/agents-layer7/pkg/config"
	"github.com/Axway/agents-layer7/pkg/models/policy"
	"github.com/Axway/agents-layer7/pkg/models/service"
	util2 "github.com/Axway/agents-layer7/pkg/util"
)

const (
	// MAJOR -
	MAJOR = "MAJOR"
	// MINOR -
	MINOR               = "MINOR"
	minorHash           = "minorHash"
	majorHash           = "majorHash"
	activePolicyVersion = "activePolicyVersion"
	xmlIndent           = "  "
)

// ServiceDetail - Sample struct representing the API definition in API gateway
type ServiceDetail struct {
	AgentDetails      map[string]string
	APIName           string
	APISpec           []byte
	APIUpdateSeverity string
	AuthPolicy        string
	Description       string
	Documentation     []byte
	Endpoint          string
	Endpoints         []apic.EndpointDefinition
	ID                string
	ResourceType      string
	ServiceAttributes map[string]string
	Title             string
}

// apiDiscovery - discovers apis
type apiDiscovery struct {
	cfg       *config.Layer7Config
	client    *client.GatewayClient
	apiCh     chan *ServiceDetail
	log       log.FieldLogger
	validator *validator
}

// newAPIDiscovery - Creates a new api apiDiscovery
func newAPIDiscovery(cfg *config.Layer7Config, c *client.GatewayClient, apiCh chan *ServiceDetail) *apiDiscovery {
	logger := log.NewFieldLogger().WithPackage("apiDiscovery").WithComponent("apiDiscovery")
	return &apiDiscovery{
		cfg:       cfg,
		client:    c,
		apiCh:     apiCh,
		log:       logger,
		validator: newValidator(),
	}
}

// Execute starts the api apiDiscovery
func (a *apiDiscovery) Execute() error {
	services, err := a.client.GetServices()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for _, item := range services.List.Item {
		wg.Add(1)

		svc := item.Resource.Service
		go func(s service.Service) {
			defer wg.Done()
			err := a.process(s)
			if err != nil {
				a.log.Errorf("failed to process service: %s", err)
			}
		}(svc)
	}

	wg.Wait()
	a.validator.SetAPIs(services.List.Item)
	return nil
}

func (a *apiDiscovery) process(svc service.Service) error {
	if svc.ServiceDetail.Enabled != "true" || isInternal(svc.ServiceDetail.Properties.Property) {
		return nil
	}

	ext := newServiceDetails(svc)

	p, err := a.client.GetActivePolicy(svc.ID)
	if err != nil {
		return fmt.Errorf("failed to get active policy: %s", err)
	}

	svcProperties := svc.ServiceDetail.Properties.Property
	activeRevision := getPolicyRevision(svcProperties)
	ext.AgentDetails[activePolicyVersion] = activeRevision

	ext.ServiceAttributes = getSvcProperties(svcProperties)

	err = a.buildService(ext, svc, p)
	if err != nil {
		a.log.Info(err)
		return nil
	}

	if ok := a.shouldPublish(ext); !ok {
		return nil
	}

	a.apiCh <- ext
	return nil
}

func (a *apiDiscovery) Status() error {
	if !a.Ready() {
		return fmt.Errorf("apiDiscovery is not running")
	}

	return nil
}

func (a *apiDiscovery) Ready() bool {
	return true
}

func (a *apiDiscovery) getAPIEndpoint(mappings service.ServiceMappings) string {
	if mappings.HTTPMapping == nil {
		return ""
	}

	if pattern, ok := mappings.HTTPMapping["UrlPattern"]; ok {
		if s, ok := pattern.(string); ok {
			s = strings.Replace(s, "*", "", -1)
			return s
		}
	}

	return ""
}

func (a *apiDiscovery) buildService(ext *ServiceDetail, svc service.Service, p *policy.PolicyItem) error {
	if isSoapAPI(svc.ServiceDetail.Properties.Property) {
		return a.processSoap(ext, svc)
	}
	return a.processService(ext, svc, p)
}

func (a *apiDiscovery) shouldPublish(ext *ServiceDetail) bool {
	minor, _ := util2.ComputeMinorHash(ext.APIName, ext.ServiceAttributes)
	// TODO: Major Hash should include the spec
	major, _ := util2.ComputeMajorHash(ext.Endpoint, ext.AgentDetails[activePolicyVersion])
	ext.AgentDetails[minorHash] = minor
	ext.AgentDetails[majorHash] = major

	savedMinor := agent.GetAttributeOnPublishedAPIByID(ext.ID, minorHash)
	if savedMinor != minor {
		ext.APIUpdateSeverity = MINOR
		a.log.Infof("minor revision update for %s", ext.APIName)
	}

	savedMajor := agent.GetAttributeOnPublishedAPIByID(ext.ID, majorHash)
	if savedMajor != major {
		ext.APIUpdateSeverity = MAJOR
		a.log.Infof("new revision update for %s", ext.APIName)
	} else if savedMinor == minor && savedMajor == major {
		a.log.Infof("no change detected for api %s", ext.APIName)
		return false
	}

	return true
}

func (a *apiDiscovery) getOASEndpoint(variables []policy.SetVariable) (string, string) {
	docHost := ""
	apiType := ""

	for _, v := range variables {
		resourceType := getResourceType(v.VariableToSet.StringValue)
		if resourceType == "" {
			continue
		}

		apiType = resourceType

		docURL, err := base64.StdEncoding.DecodeString(v.Base64Expression.StringValue)
		if err != nil {
			a.log.Error("failed to decode host: %s", err)
			continue
		}
		docHost = string(docURL)

		break
	}

	return apiType, docHost
}

func (a *apiDiscovery) processSoap(ext *ServiceDetail, svc service.Service) error {
	endpoint := a.getAPIEndpoint(svc.ServiceDetail.ServiceMappings)
	if endpoint == "" {
		return fmt.Errorf("unable to find proxy endpoint for %s", svc.ServiceDetail.Name)
	}

	content := a.getWsdl(svc.Resources.ResourceSet)
	if content == nil {
		return fmt.Errorf("unable to find wsdl spec for %s", svc.ServiceDetail.Name)
	}

	ext.APISpec = content
	ext.ResourceType = apic.Wsdl
	ext.Endpoint = endpoint
	url := a.cfg.Host + endpoint
	ep := util2.CreateEndpoint(url)
	ext.Endpoints = append(ext.Endpoints, ep)
	return nil
}

func (a *apiDiscovery) processService(ext *ServiceDetail, svc service.Service, p *policy.PolicyItem) error {
	endpoint := a.getAPIEndpoint(svc.ServiceDetail.ServiceMappings)
	if endpoint == "" {
		return fmt.Errorf("unable to find proxy endpoint for %s", svc.ServiceDetail.Name)
	}
	ext.Endpoint = endpoint
	url := a.cfg.Host + endpoint
	ep := util2.CreateEndpoint(url)
	ext.Endpoints = append(ext.Endpoints, ep)

	resourceType, docHost := a.getOASEndpoint(p.Policy.All.SetVariable)

	if resourceType != "" && docHost != "" {
		// OAS Spec
		return a.processOAS(ext, resourceType, docHost)
	}

	// Unstructured
	unstructuredSvc, err := xml.MarshalIndent(svc, "", xmlIndent)
	if err != nil {
		return err
	}

	ext.ResourceType = apic.Unstructured
	ext.APISpec = unstructuredSvc
	return nil
}

func (a *apiDiscovery) processOAS(ext *ServiceDetail, resourceType, docHost string) error {
	spec, err := a.client.GetSpec(docHost)
	if err != nil {
		return fmt.Errorf("failed to get spec: %s", err)
	}

	ext.APISpec = spec
	ext.ResourceType = resourceType

	return nil
}

func (a *apiDiscovery) getWsdl(resources []service.ResourceSetElement) []byte {
	for _, item := range resources {
		if item.Tag == "wsdl" {
			for _, res := range item.Resource {
				if res.Type == "wsdl" {
					return []byte(res.Content)
				}
			}
		}
	}

	return nil
}
