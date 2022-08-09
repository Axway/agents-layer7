package config

import (
	"fmt"
	"time"

	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	"github.com/Axway/agent-sdk/pkg/cmd/properties"
	corecfg "github.com/Axway/agent-sdk/pkg/config"
)

const (
	pathUsername     = "layer7.username"
	pathPass         = "layer7.password"
	pathHost         = "layer7.host"
	pathAPI          = "layer7.api"
	pathPollInterval = "layer7.pollInterval"
)

// AgentConfig -
type AgentConfig struct {
	CentralCfg corecfg.CentralConfig `config:"central"`
	Layer7Cfg  *Layer7Config         `config:"layer7"`
}

// Layer7Config -
type Layer7Config struct {
	Username     string        `config:"username"`
	Password     string        `config:"password"`
	Host         string        `config:"host"`
	API          string        `config:"api"`
	PollInterval time.Duration `config:"pollInterval"`
}

// ValidateCfg -
func (c *Layer7Config) ValidateCfg() (err error) {
	if c.Username == "" {
		return fmt.Errorf("username not provided")
	}

	if c.Password == "" {
		return fmt.Errorf("password not provided")
	}

	if c.Host == "" {
		return fmt.Errorf("host not provided")
	}

	if c.API == "" {
		return fmt.Errorf("api endpoint not provided")
	}

	return
}

// ApplyResources -
func (c *Layer7Config) ApplyResources(agentResource *v1.ResourceInstance) error {
	return nil
}

// AddProperties - adds config needed for apigee client
func AddProperties(rootProps properties.Properties) {
	rootProps.AddStringProperty(pathUsername, "", "Layer7 username")
	rootProps.AddStringProperty(pathPass, "", "Layer7 password")
	rootProps.AddStringProperty(pathHost, "", "Layer7 hostname")
	rootProps.AddStringProperty(pathAPI, "", "Layer7 REST API endpoint")
	rootProps.AddDurationProperty(pathPollInterval, 1*time.Minute, "The time interval to check for Layer7 services")
}

// ParseConfig - parse the config on startup
func ParseConfig(rootProps properties.Properties) *Layer7Config {
	return &Layer7Config{
		Username:     rootProps.StringPropertyValue(pathUsername),
		Password:     rootProps.StringPropertyValue(pathPass),
		Host:         rootProps.StringPropertyValue(pathHost),
		API:          rootProps.StringPropertyValue(pathAPI),
		PollInterval: rootProps.DurationPropertyValue(pathPollInterval),
	}
}
