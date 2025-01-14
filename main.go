package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/implementation"
	"github.com/Tether-Payments/tamarin-mux/tamarin"
)

var port int

func main() {
	server := tamarin.NewServer().
		WithEndpoint(tamarin.NewEndpoint("/fancyping", http.MethodPost).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.EndpointPingPOST)).
		WithHandleFuncs("/ping", http.MethodGet, implementation.PingGET).
		WithHandleFuncs("/ping", http.MethodPost, implementation.PingPOST).
		WithEndpoint(tamarin.NewEndpoint("/failjson", http.MethodPost).WithHandlers(implementation.FailIfNoBody, implementation.ShowBody))

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("Now listening on %s with handlers for", addr)
	for _, name := range server.HandlerNames() {
		log.Print(name)
	}
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
}

func init() {
	flag.IntVar(&port, "port", 12345, "The port on which to listen")

	flag.Parse()
}
