package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/examples/implementation"
	"github.com/Tether-Payments/tamarin-mux/tamarin"
)

var port int

func main() {
	handler := tamarin.NewHandler(true).
		WithEndpoint(tamarin.NewEndpoint("/fancyping", http.MethodPost).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.EndpointPingPOST)).
		WithHandleFuncs("/ping", http.MethodGet, implementation.PingGET).
		WithHandleFuncs("/ping", http.MethodPost, implementation.PingPOST).
		WithEndpoint(tamarin.NewEndpoint("/wallet/{}", http.MethodGet).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.PrintURLWithElements)).
		WithEndpoint(tamarin.NewEndpoint("/wallet/{}/something/{}/else", http.MethodGet).WithHandlers(implementation.PrintURLWithElements)).
		WithEndpoint(tamarin.NewEndpoint("/failjson", http.MethodPost).WithHandlers(implementation.FailIfNoBody, implementation.ShowBody)).
		WithEndpoint(tamarin.NewEndpoint("/content/{*}", http.MethodPost).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.StaticSiteHandler)).
		WithEndpoint(tamarin.NewEndpoint("/content/{*}", http.MethodGet).WithHandlers(implementation.StaticSiteHandler))

	addr := fmt.Sprintf("0.0.0.0:%d", port)

	log.Printf("Now listening on %s with handlers for", addr)
	for _, name := range handler.HandlerNames() {
		log.Print(name)
	}
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}
}

func init() {
	flag.IntVar(&port, "port", 12345, "The port on which to listen")

	flag.Parse()
}
