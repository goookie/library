package consul

import (
	consulWatch "github.com/hashicorp/consul/api/watch"
)

// ClientInterface defines the Client of Consul
// Registration/Registration service to Consul
// Listening to the service in Consul
type ClientInterface interface {
	Wait()                                           // Wait for a specific service to come online
	Register() error                                 // Registration service to Consul
	DeRegister() error                               // DeRegister service to Consul
	Watch(serverType string) <-chan AvailableServers // Listening to the service in Consul
}

// AvailableServers defines available online services
type AvailableServers struct {
	ServerType string
	Servers    []string
}

// RegistryConfig is service registry config
type RegistryConfig struct {
	ID         string   // service id
	ServerType string   // service type
	IP         string   // service addr
	Port       int      // service port
	Tags       []string // service Tags
}

// DiscoveryConfig is service watcher config
type DiscoveryConfig struct {
	ServerType string   // target service type
	Tags       []string // target service tags
	Min        int      // minimum of available in wait
}

func (discoveryConfig *DiscoveryConfig) getWatcher() (*watcher, error) {
	w := &watcher{
		serverType: discoveryConfig.ServerType,
		noticeChan: make(chan AvailableServers, 100),
	}

	// build plan
	params := make(map[string]interface{})
	params["type"] = "service"
	params["service"] = discoveryConfig.ServerType
	params["tag"] = discoveryConfig.Tags
	plan, err := consulWatch.Parse(params)
	if err != nil {
		return nil, err
	}
	plan.Handler = w.handler

	w.plan = plan
	return w, nil
}
