package consul

import (
	consulAPI "github.com/hashicorp/consul/api"
)

// Client defines the consul client
type Client struct {
	config       *consulAPI.Config // consul config
	consulClient *consulAPI.Client // consul Client

	// service registration related
	registryConfig *RegistryConfig

	// service watcher related
	discoveryConfigs map[string]*DiscoveryConfig
}

type ClientOption func(*Client) error

func NewClient(opts ...ClientOption) (*Client, error) {
	c := new(Client)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func AddrOption(config *consulAPI.Config) ClientOption {
	return func(client *Client) error {
		// service registry
		c, err := consulAPI.NewClient(config)
		if err != nil {
			return err
		}

		client.config = config
		client.consulClient = c
		return nil
	}
}

func ServiceRegistryOption(registryConfig *RegistryConfig) ClientOption {
	return func(client *Client) error {
		client.registryConfig = registryConfig
		return nil
	}
}

func ServiceDiscoveryOption(discoveryConfigs ...*DiscoveryConfig) ClientOption {
	return func(client *Client) error {
		// service watcher channel
		configs := make(map[string]*DiscoveryConfig, 0)

		// service watcher
		for _, sdConfig := range discoveryConfigs {
			// add to service watcher configs
			configs[sdConfig.ServerType] = sdConfig
		}

		client.discoveryConfigs = configs
		return nil
	}
}
