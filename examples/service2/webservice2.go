// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu
package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/werbenhu/registry/register"
)

var (
	mu        sync.Mutex
	WebGroup  = "webservice-group"
	ServiceId = "webserver2"
)

// login function
func login(c *gin.Context) {
	userid := c.Query("userid")
	c.JSON(http.StatusOK, map[string]any{
		"msg":    "success from:" + ServiceId,
		"userid": userid,
	})
}

func main() {
	var err error
	registries := "172.16.3.3:7370"
	bind := ":8371"
	bindAdvertise := "172.16.3.3:8371"
	addr := "172.16.3.3:8001"

	// New() create a register object
	// id:
	//    service id
	// bind:
	//    The address used to register the service to registry server.
	//    If there is a firewall, please remember that the port needs to open both tcp and udp.
	// advertise:
	//    The address that the service will advertise to registry server.
	//    Can be used for basic NAT traversal where both the internal ip:port and external ip:port are known.
	// registries:
	//    The addresses of the registry servers, if there are more than one, separate them with commas,
	//    such as "192.168.1.101:7370,192.168.1.102:7370"
	// group:
	//    Group name the current service belongs to.
	// addr:
	//    The address currently provided by this service to the client,
	//    for example, the current service is an http server,
	//    that is the address 172.16.3.3:80 that http listens to.
	reg := register.New(ServiceId, bind, bindAdvertise, registries, WebGroup, addr)
	err = reg.Start()
	if err != nil {
		panic(err)
	}

	// start web service
	r := gin.Default()
	r.GET("/login", login)
	r.Run(addr)
}
