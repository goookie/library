package consul

import (
	"log"
)

// Watch listening to the service in Consul
func (client *Client) Watch(serverType string) <-chan AvailableServers {
	if len(client.discoveryConfigs) == 0 {
		return nil
	}

	sdConfig, isExist := client.discoveryConfigs[serverType]
	if isExist {
		client.once.Do(func() {
			for _, sdConfig := range client.discoveryConfigs {
				go func(sdConfig *DiscoveryConfig) {
					if err := sdConfig.plan.Run(client.consulAddr); err != nil {
						log.Printf("Consul Watch Err: %+v\n", err)
					}
				}(sdConfig)
			}
		})
		return sdConfig.watchChan
	}

	return nil
}
