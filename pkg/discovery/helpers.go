package discovery

import (
	"strings"

	"git.ecd.axway.org/tjohnson/layer7/pkg/models/service"
	"github.com/Axway/agent-sdk/pkg/apic"
)

func getPolicyRevision(props []service.Property) string {
	for _, prop := range props {
		if prop.Key == "policyRevision" {
			return *prop.LongValue
		}
	}
	return ""
}

func isSoapAPI(properties []service.Property) bool {
	for _, prop := range properties {
		if prop.Key == "soap" && *prop.BooleanValue == "true" {
			return true
		}
	}

	return false
}

func isInternal(properties []service.Property) bool {
	for _, prop := range properties {
		if prop.Key == "internal" && *prop.BooleanValue == "true" {
			return true
		}
	}

	return false
}

func getSvcProperties(props []service.Property) map[string]string {
	attrs := make(map[string]string)
	for _, prop := range props {
		if ok := strings.Contains(prop.Key, "property."); ok {
			split := strings.Split(prop.Key, ".")
			if len(split) == 2 {
				attrs[split[1]] = *prop.StringValue
			}
		}
	}

	return attrs
}

func getResourceType(urlVar string) string {
	switch urlVar {
	case "openapi.docUrl":
		return apic.Oas3
	case "swagger.docUrl":
		return apic.Oas2
	default:
		return ""
	}
}

func newServiceDetails(svc service.Service) *ServiceDetail {
	return &ServiceDetail{
		APIName:           svc.ServiceDetail.Name,
		APISpec:           nil,
		APIUpdateSeverity: "",
		AgentDetails:      map[string]string{},
		ServiceAttributes: map[string]string{},
		ID:                svc.ID,
		Title:             svc.ServiceDetail.Name,
		ResourceType:      "",
	}
}
