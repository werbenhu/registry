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
	id := flag.String("id", "node1", "")
	addr := flag.String("addr", "172.16.3.3:7370", "")
	advertise := flag.String("advertise", "172.16.3.3:7370", "")
	members := flag.String("members", "", "")
	service := flag.String("service", "", "")
	port := flag.String("api-port", "8080", "")

	flag.Parse()
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
		srouter.OptMembers(*members),
		srouter.OptApiPort(*port),
		srouter.OptService(*service),
	})

	err := router.Serve()
	if err != nil {
		log.Fatal(err)
	}

	<-done
	router.Close()
}
