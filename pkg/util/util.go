package util

import (
	"net/url"
	"strconv"

	"github.com/Axway/agent-sdk/pkg/apic"
	"github.com/Axway/agent-sdk/pkg/util"
)

func ComputeMinorHash(name string, attrs map[string]string) (string, error) {
	m := make(map[string]string)
	m["serviceName"] = name
	m = util.MergeMapStringString(m, attrs)
	hashed, err := util.ComputeHash(m)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(hashed, 10), nil
}

func ComputeMajorHash(ep, version string) (string, error) {
	m := map[string]string{
		"endpoint": ep,
		"version":  version,
	}
	hashed, err := util.ComputeHash(m)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(hashed, 10), nil
}

func CreateEndpoint(host string) apic.EndpointDefinition {
	parsed, _ := url.Parse(host)
	port, _ := strconv.Atoi(parsed.Port())

	return apic.EndpointDefinition{
		Host:     parsed.Hostname(),
		Port:     int32(port),
		Protocol: parsed.Scheme,
		BasePath: parsed.Path,
	}
}
