package consul

import (
	"log"
)

// Watch listening to the service in Consul
func (client *Client) Watch(serverType string) <-chan AvailableServers {
	sdConfig, isExist := client.discoveryConfigs[serverType]
	if !isExist {
		return nil
	}

	w, err := sdConfig.getWatcher()
	if err != nil {
		log.Fatalf("Consul Watch Err: %+v\n", err)
	}

	go func(w *watcher) {
		if err := w.plan.Run(client.config.Address); err != nil {
			log.Printf("Consul Watch Err: %+v\n", err)
		}
	}(w)

	return w.noticeChan
}
