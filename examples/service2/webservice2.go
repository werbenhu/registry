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

// login handles the login request.
func login(c *gin.Context) {
	userid := c.Query("userid")
	c.JSON(http.StatusOK, map[string]any{
		"msg":    "success from:" + ServiceId,
		"userid": userid,
	})
}

func main() {
	var err error

	// Configuration for registry registration
	registries := "172.16.3.3:7370"
	bind := ":8371"
	bindAdvertise := "172.16.3.3:8371"
	addr := "172.16.3.3:8001"

	// Create a new register object.
	// id: The service id.
	// bind: The address used to register the service to the registry server.
	// advertise: The address that the service will advertise to the registry server.
	// registries: The addresses of the registry servers, separated by commas if there are more than one.
	// group: The group name the current service belongs to.
	// addr: The address currently provided by this service to the client.
	reg := register.New(ServiceId, bind, bindAdvertise, registries, WebGroup, addr)

	// Start the registry.
	err = reg.Start()
	if err != nil {
		panic(err)
	}

	// Start the web service.
	r := gin.Default()
	r.GET("/login", login)
	r.Run(addr)
}
