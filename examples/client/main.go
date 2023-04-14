package main

import (
	"fmt"
	"log"

	"github.com/werbenhu/registry/client"
)

func main() {
	// 路由服务器中任意选择一个都可以
	routerService := "172.16.3.3:9001"
	group := "webservice-group"

	client, err := client.NewRpcClient(routerService)
	if err != nil {
		panic(err)
	}

	// 根据用户ID使用一致性哈希分配服务
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("user-id-%d", i)
		service, err := client.Match(group, key)

		if err != nil {
			log.Printf("[ERROR] match key%s err:%s\n", key, err)
			continue
		}
		log.Printf("[INFO] match key:%s, serviceId:%s, serviceAddr:%s\n", key, service.Id, service.Addr)
	}

	// 获取webservice-group组所有的服务
	allService, err := client.Members(group)
	if err != nil {
		log.Printf("[ERROR] get all service err:%s\n", err)
	}
	log.Printf("[INFO] all service:%+v\n", allService)
}
