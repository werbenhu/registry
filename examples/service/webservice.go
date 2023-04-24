// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu
package main

import (
	"flag"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/werbenhu/registry/register"
)

var (
	mu sync.Mutex
)

// login handles the login request.
func login(c *gin.Context) {
	userid := c.Query("userid")
	c.JSON(http.StatusOK, map[string]any{
		"msg":    "success",
		"userid": userid,
	})
}

func main() {
	var err error

	// Configuration for registry registration
	registries := flag.String("registries", "", "Registry server addresses, it can be empty, and multiples are separated by commas.")
	id := flag.String("id", "", "The service id")
	addr := flag.String("addr", ":8000", "The address used for service discovery (default \":9800\").")
	advertise := flag.String("advertise", "", "The address will advertise to client for service discovery (default \":9800\").")
	bind := flag.String("bind", ":7370", "The address used to register the service (default \":7370\").")
	group := flag.String("group", "webservice-group", "The group name the current service belongs to")
	flag.Parse()

	// If the Service ID not set, a random one will be generated
	if *id == "" {
		*id = xid.New().String()
	}

	// Create a new register object.
	// id: The service id.
	// bind: The address used to register the service to the registry server.
	// advertise: The address that the service will advertise to the registry server.
	// registries: The addresses of the registry servers, separated by commas if there are more than one.
	// group: The group name the current service belongs to.
	// addr: The address currently provided by this service to the client.
	reg := register.New(*id, *bind, *advertise, *registries, *group, *addr)

	// Start the registry.
	err = reg.Start()
	if err != nil {
		panic(err)
	}

	// Start the web service.
	r := gin.Default()
	r.GET("/login", login)
	r.Run(*addr)
}
