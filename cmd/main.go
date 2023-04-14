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

	"github.com/werbenhu/registry"
)

func main() {
	id := flag.String("id", "", "The service id, cannot be empty")
	bind := flag.String("bind", ":7370", "The address used to register the service (default \":7370\").")
	bindAdvertise := flag.String("bind-advertise", ":7370", "The address will advertise to other services (default \":7370\").")
	registries := flag.String("registries", "", "Registry server addresses, it can be empty, and multiples are separated by commas.")
	addr := flag.String("addr", ":9800", "The address used for service discovery (default \":9800\").")
	advertise := flag.String("advertise", "", "The address will advertise to client for service discover (default \":9800\").")

	flag.Parse()
	if *id == "" {
		log.Fatal(registry.ErrMemberIdEmpty)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	r := registry.New([]registry.IOption{
		registry.OptId(*id),
		registry.OptBind(*bind),
		registry.OptBindAdvertise(*bindAdvertise),
		registry.OptAddr(*addr),
		registry.OptAdvertise(*advertise),
		registry.OptRegistries(*registries),
	})

	err := r.Serve()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] registry server start finished.\n")
	<-done
	r.Close()
}
