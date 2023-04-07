// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu
package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/werbenhu/srouter/register"
)

var (
	mu        sync.Mutex
	WebGroup  = "webservice-group"
	ServiceId = "webserver2"
)

// 登录接口
func login(c *gin.Context) {
	userid := c.Query("userid")
	c.JSON(http.StatusOK, map[string]any{
		"msg":    "success from:" + ServiceId,
		"userid": userid,
	})
}

func main() {
	var err error

	// 路由服务器的地址，服务自动发现的时候，需要注册到路由服务器去，
	// 多个用逗号隔开，这里只需要一个，路由服务器中任意选择一个都可以
	routers := "172.16.3.3:7370"

	// addr是当前本服务需要和路由服务器通信的地址，
	// 如果有防火墙，请记得端口需要同时打开tcp和udp
	addr := ":8371"

	// 对外公布的服务发现通信的地址，需要这个参数涉及到网关有端口映射的时候，
	// 比如docker，内部监听的地址是127.0.0.1:80, 映射到对外则是公网IP:8000
	advertise := "172.16.3.3:8371"

	// 当前本服务提供服务的地址，比如当前服务是http服务器，那就是http监听的那个地址172.16.3.3:8000
	service := "172.16.3.3:8001"

	// 将当前服务注册到路由服务器
	reg := register.New(ServiceId, addr, advertise, routers, WebGroup, service)
	err = reg.Run()
	if err != nil {
		panic(err)
	}

	// 启动http服务
	r := gin.Default()
	r.GET("/login", login)
	r.Run(service)
}
