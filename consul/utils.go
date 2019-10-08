package consul

import (
	"fmt"

	consulAPI "github.com/hashicorp/consul/api"
	consulWatch "github.com/hashicorp/consul/api/watch"
)

type watcher struct {
	serverType string
	plan       *consulWatch.Plan
	noticeChan chan AvailableServers
}

func (d *watcher) handler(index uint64, raw interface{}) {
	if raw == nil {
		return
	}
	if entries, ok := raw.([]*consulAPI.ServiceEntry); ok {
		var servers []string
		for _, entry := range entries {
			// healthy check fail, continue anyway
			if entry.Checks.AggregatedStatus() != consulAPI.HealthPassing {
				continue
			}
			servers = append(servers, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
		}
		d.noticeChan <- AvailableServers{
			ServerType: d.serverType,
			Servers:    servers,
		}
	}
}
