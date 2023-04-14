package main

import (
	"fmt"
	"log"

	"github.com/werbenhu/registry/client"
)

func main() {
	// registry server address
	registryAddr := "172.16.3.3:9801"
	group := "webservice-group"

	client, err := client.NewRpcClient(registryAddr)
	if err != nil {
		panic(err)
	}

	// assign services to 1000 users with consistent hashing algorithm
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user-id-%d", i)
		service, err := client.Match(group, key)

		if err != nil {
			log.Printf("[ERROR] match key%s err:%s\n", key, err)
			continue
		}
		log.Printf("[INFO] match key:%s, serviceId:%s, serviceAddr:%s\n", key, service.Id, service.Addr)
	}

	// get all service of the group
	allService, err := client.Members(group)
	if err != nil {
		log.Printf("[ERROR] get all service err:%s\n", err)
	}
	log.Printf("[INFO] all service:%+v\n", allService)
}
