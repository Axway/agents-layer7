package cmd

import (
	corecmd "github.com/Axway/agent-sdk/pkg/cmd"
	corecfg "github.com/Axway/agent-sdk/pkg/config"

	"git.ecd.axway.org/tjohnson/layer7/pkg/config"
	"git.ecd.axway.org/tjohnson/layer7/pkg/discovery"
)

// RootCmd - Agent root command
var RootCmd corecmd.AgentRootCmd
var discoveryAgent *discovery.Agent

func init() {
	RootCmd = corecmd.NewRootCmd(
		"layer7_discovery_agent",
		"Layer7 Discovery Agent",
		initConfig,
		run,
		corecfg.DiscoveryAgent,
	)

	rootProps := RootCmd.GetProperties()
	config.AddProperties(rootProps)
}

func run() error {
	return discoveryAgent.Run()
}

func initConfig(centralConfig corecfg.CentralConfig) (interface{}, error) {
	rootProps := RootCmd.GetProperties()
	l7Cfg := config.ParseConfig(rootProps)

	agentConfig := config.AgentConfig{
		CentralCfg: centralConfig,
		Layer7Cfg:  l7Cfg,
	}

	discoveryAgent = discovery.NewAgent(agentConfig)

	return agentConfig, nil
}
