package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	discovery "github.com/andy-trimble/go.discovery"
)

func main() {
	d := discovery.Discovery{}
	err := d.Start("server")
	if err != nil {
		log.Fatal(err)
	}
	defer d.Shutdown()

	go func() {
		for {
			actor := <-d.Discovered
			log.Printf("Discovered %+v", actor)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		d.Shutdown()
		os.Exit(0)
	}()

	for {
		err := <-d.Err
		log.Printf("%+v", err)
	}
}
