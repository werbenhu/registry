package main

import (
	"fmt"
	"log"

	"github.com/werbenhu/registry/client"
)

func main() {
	// Address of the registry server.
	registryAddr := "172.16.3.3:9801"

	// Group name of the services to be matched.
	group := "webservice-group"

	// Create a new RPC client.
	client, err := client.NewRpcClient(registryAddr)
	if err != nil {
		panic(err)
	}

	// Assign services to 1000 users with consistent hashing algorithm.
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user-id-%d", i)

		// Find the service assigned to the given key.
		service, err := client.Match(group, key)
		if err != nil {
			log.Printf("[ERROR] Failed to match key %s: %s\n", key, err)
			continue
		}
		log.Printf("[INFO] Matched key: %s, Service ID: %s, Service Address: %s\n", key, service.Id, service.Addr)
	}

	// Get all services of the group.
	allService, err := client.Members(group)
	if err != nil {
		log.Printf("[ERROR] Failed to get all services: %s\n", err)
	}
	log.Printf("[INFO] All services: %+v\n", allService)
}
