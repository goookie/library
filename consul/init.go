package consul

import (
	"sync"

	consulAPI "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/api/watch"
)

// Client defines the consul client
type Client struct {
	consulAddr string // consul address（127.0.0.1:8500）

	// service registration related
	registryConfig *RegistryConfig
	consulClient   *consulAPI.Client // consul Client

	// service discovery related
	once             sync.Once
	discoveryConfigs map[string]*DiscoveryConfig
	//watchChan        chan AvailableServers
	//watchChannels    map[string]chan AvailableServers
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

func AddrOption(addr string) ClientOption {
	return func(client *Client) error {
		// service registry
		c, err := consulAPI.NewClient(&consulAPI.Config{Address: addr})
		if err != nil {
			return err
		}

		client.consulAddr = addr
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
		// service discovery channel
		configs := make(map[string]*DiscoveryConfig, 0)

		// service discovery
		for _, sdConfig := range discoveryConfigs {
			// build watch chan
			watchChan := make(chan AvailableServers, 100)

			// build plan
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = sdConfig.ServerType
			params["tag"] = sdConfig.Tags
			plan, err := consulWatch.Parse(params)
			if err != nil {
				return err
			}
			plan.Handler = sdConfig.handler

			// bind plan to DiscoveryConfig
			sdConfig.watchChan = watchChan
			sdConfig.plan = plan

			// add to service discovery configs
			configs[sdConfig.ServerType] = sdConfig
		}

		client.discoveryConfigs = configs
		//client.watchChan = watchChan
		client.once = sync.Once{}
		return nil
	}
}
