package consul

import (
	"log"
)

type serviceDetail struct {
	min     int
	current int
}

func (serviceDetail *serviceDetail) available() bool {
	return serviceDetail.current >= serviceDetail.min
}

// Wait wait for required service
func (client *Client) Wait() {
	// new services count
	services := make(map[string]*serviceDetail)
	for _, v := range client.discoveryConfigs {
		services[v.ServerType] = &serviceDetail{
			min: v.Min,
		}
	}

	// check first
	if allAvailable(services) {
		return
	}

	// build notice chan
	noticeChan := make(chan AvailableServers, 100)
	for _, sdConfig := range client.discoveryConfigs {
		w, err := sdConfig.getWatcher()
		if err != nil {
			log.Fatalf("consul wait err: %v\n", err)
		}

		go func(w *watcher) {
			if err := w.plan.Run(client.config.Address); err != nil {
				log.Printf("Consul Watch Err: %+v\n", err)
			}
		}(w)

		go func(noticeChan, subChan chan AvailableServers) {
			for msg := range subChan {
				noticeChan <- msg
			}
		}(noticeChan, w.noticeChan)
	}

	for msg := range noticeChan {
		detail := services[msg.ServerType]
		detail.current = len(msg.Servers)
		if allAvailable(services) {
			return
		}
	}
}

// check is all type services available
func allAvailable(servers map[string]*serviceDetail) bool {
	for _, server := range servers {
		if !server.available() {
			return false
		}
	}
	return true
}
