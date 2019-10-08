package consul

import (
	"fmt"
	"net/http"

	consulAPI "github.com/hashicorp/consul/api"
)

// Client defines the consul client
type Client struct {
	config       *consulAPI.Config // consul config
	consulClient *consulAPI.Client // consul Client

	// service registration related
	registryConfig    *RegistryConfig
	consulCheckServer *http.Server
	consulCheckPort   int

	// service watcher related
	discoveryConfigs map[string]*DiscoveryConfig
}

type ClientOption func(*Client) error

func NewClient(config *consulAPI.Config, opts ...ClientOption) (*Client, error) {
	c := new(Client)
	// service registry
	client, err := consulAPI.NewClient(config)
	if err != nil {
		return nil, err
	}

	c.config = config
	c.consulClient = client

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func ServiceRegistryOption(checkPort int, registryConfig *RegistryConfig) ClientOption {
	return func(client *Client) error {
		mux := http.NewServeMux()
		mux.HandleFunc("/check", func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("ok"))
		})

		client.registryConfig = registryConfig
		client.consulCheckPort = checkPort
		client.consulCheckServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", checkPort),
			Handler: mux,
		}

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
