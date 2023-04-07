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
	addr := flag.String("addr", ":7370", "")
	advertise := flag.String("advertise", ":7370", "")
	routers := flag.String("routers", "", "")
	apiAddr := flag.String("api-addr", ":8080", "")
	apiAdvertise := flag.String("api-advertise", "", "")

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
		srouter.OptRouters(*routers),
		srouter.OptApiAddr(*apiAddr),
		srouter.OptService(*apiAdvertise),
	})

	err := router.Serve()
	if err != nil {
		log.Fatal(err)
	}

	<-done
	router.Close()
}
