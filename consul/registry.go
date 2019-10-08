package consul

import (
	"context"
	"fmt"
	"log"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
)

var ErrAbsentServiceRegisterConfig = fmt.Errorf("service register config is absent")

// Register implements registry Client interface
func (client *Client) Register() error {
	if client.registryConfig == nil {
		return ErrAbsentServiceRegisterConfig
	}

	go func() {
		if err := client.consulCheckServer.ListenAndServe(); err != nil {
			log.Printf("cnosul check server start fail : %v\n", err)
		}
	}()

	registration := new(consulAPI.AgentServiceRegistration)
	registration.ID = client.registryConfig.ID
	registration.Name = client.registryConfig.ServerType
	registration.Tags = client.registryConfig.Tags
	registration.Address = client.registryConfig.IP
	registration.Port = client.registryConfig.Port
	registration.Check = &consulAPI.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, client.consulCheckPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s", // del this service in 15s after check fail
	}

	return client.consulClient.Agent().ServiceRegister(registration)
}

func (client *Client) DeRegister() error {
	if client.registryConfig == nil {
		return ErrAbsentServiceRegisterConfig
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client.consulCheckServer.Shutdown(ctx)
	defer cancel()

	return client.consulClient.Agent().ServiceDeregister(client.registryConfig.ID)
}
