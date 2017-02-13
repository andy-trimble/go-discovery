package main

import (
	"bufio"
	"log"
	"os"

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

	go func() {
		for {
			err := <-d.Err
			log.Printf("%+v", err)
		}
	}()

	log.Println("Press return to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
