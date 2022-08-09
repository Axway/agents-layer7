package discovery

import (
	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/apic"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/Axway/agent-sdk/pkg/util/log"
)

type publisher struct {
	publishAPI agent.PublishAPIFunc
	stop       chan interface{}
	apiCh      chan *ServiceDetail
	log        log.FieldLogger
}

func newPublisher() *publisher {
	logger := log.NewFieldLogger().WithPackage("discovery").WithComponent("publisher")
	return &publisher{
		publishAPI: agent.PublishAPI,
		stop:       make(chan interface{}),
		apiCh:      make(chan *ServiceDetail),
		log:        logger,
	}
}

// Execute starts the api publisher channel
func (p *publisher) Execute() error {
	for {
		select {
		case <-p.stop:
			return nil
		case api := <-p.apiCh:
			err := p.publish(api)
			if err != nil {
				p.log.Errorf("failed to published api %s: %s", api.APIName, err)
			}
		}
	}
}

// Status -
func (p *publisher) Status() error {
	return nil
}

// Ready -
func (p *publisher) Ready() bool {
	return true
}

// publish the API to Amplify Central.
func (p *publisher) publish(serviceDetail *ServiceDetail) error {
	p.log.Infof("publishing to Amplify Central")

	serviceBody, err := BuildServiceBody(serviceDetail)
	if err != nil {
		return err
	}
	err = p.publishAPI(serviceBody)
	if err != nil {
		return err
	}
	p.log.Infof("published API to Amplify Central")
	return nil
}

// BuildServiceBody - creates the service definition
func BuildServiceBody(service *ServiceDetail) (apic.ServiceBody, error) {
	return apic.NewServiceBodyBuilder().
		SetAPISpec(service.APISpec).
		SetAPIUpdateSeverity(service.APIUpdateSeverity).
		SetAuthPolicy(service.AuthPolicy).
		SetDescription(service.Description).
		SetDocumentation(service.Documentation).
		SetID(service.ID).
		SetResourceType(service.ResourceType).
		SetServiceAgentDetails(util.MapStringStringToMapStringInterface(service.AgentDetails)).
		SetServiceAttribute(service.ServiceAttributes).
		SetTitle(service.Title).
		SetServiceEndpoints(service.Endpoints).
		Build()
}
