package discovery

import (
	"git.ecd.axway.org/tjohnson/layer7/pkg/client"
	"git.ecd.axway.org/tjohnson/layer7/pkg/config"
	"github.com/Axway/agent-sdk/pkg/agent"
	"github.com/Axway/agent-sdk/pkg/jobs"
	"github.com/Axway/agent-sdk/pkg/util/log"
)

// NewAgent creates an agent for API Discovery
func NewAgent(cfg config.AgentConfig) *Agent {
	logger := log.NewFieldLogger().WithPackage("apiDiscovery").WithComponent("agent")
	return &Agent{
		Cfg:  cfg,
		log:  logger,
		stop: make(chan struct{}),
	}
}

// Agent struct for running the agent
type Agent struct {
	Cfg  config.AgentConfig
	stop chan struct{}
	log  log.FieldLogger
}

// Run starts the agent
func (a *Agent) Run() error {
	go func() {
		err := a.register()
		if err != nil {
			a.log.Errorf("failed to register: %s", err)
			a.Stop()
		}
	}()

	<-a.stop

	return nil
}

// Stop stops the agent
func (a *Agent) Stop() {
	a.stop <- struct{}{}
}

// register jobs for apiDiscovery
func (a *Agent) register() error {
	l7Cfg := a.Cfg.Layer7Cfg

	p := newPublisher()
	c := client.NewGatewayClient(l7Cfg.Host+l7Cfg.API, l7Cfg.Username, l7Cfg.Password)
	disc := newAPIDiscovery(l7Cfg, c, p.apiCh)

	_, err := jobs.RegisterIntervalJobWithName(disc, l7Cfg.PollInterval, "api-discovery")
	if err != nil {
		return err
	}

	_, err = jobs.RegisterChannelJobWithName(p, p.stop, "api-publisher")
	if err != nil {
		return err
	}

	agent.RegisterAPIValidator(disc.validator.Validate)

	return nil
}
