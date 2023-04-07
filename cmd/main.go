// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/werbenhu/srouter"
)

func main() {
	id := flag.String("id", "", "服务ID，不能为空")
	addr := flag.String("addr", ":7370", "服务发现通信的地址")
	advertise := flag.String("advertise", ":7370", "对外公布的服务发现通信的地址")
	routers := flag.String("routers", "", "路由服务器地址，如果是第一个可以为空，多个用逗号隔开")
	apiAddr := flag.String("api-addr", ":8080", "查询服务器的地址")
	service := flag.String("service", "", "对外公布的查询服务器的地址")

	flag.Parse()
	if *id == "" {
		log.Fatal(srouter.ErrMemberIdEmpty)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	router := srouter.New([]srouter.IOption{
		srouter.OptId(*id),
		srouter.OptAddr(*addr),
		srouter.OptAdvertise(*advertise),
		srouter.OptRouters(*routers),
		srouter.OptApiAddr(*apiAddr),
		srouter.OptService(*service),
	})

	err := router.Serve()
	if err != nil {
		log.Fatal(err)
	}

	<-done
	router.Close()
}
