package servicediscoveryconsul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	port      = 3003
	serviceID = "cart-serive"
)

func RegisterServie() {
	config := consulapi.DefaultConfig()

	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Println(err.Error())
		return
	}

	addr := "localhost"

	log.Println(addr)

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "cart-server",
		Port:    port,
		Address: addr,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", addr, port, serviceID),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	log.Printf(fmt.Sprintf("%s:%d/%s", addr, port, serviceID))

	regiErr := consul.Agent().ServiceRegister(registration)
	if regiErr != nil {
		log.Printf("failed to register service: %s:%v", addr, port)
	} else {
		log.Printf("successfully register service %s:%v", addr, port)
	}
}
