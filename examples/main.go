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
		WithPostEndpoint(tamarin.NewEndpoint("/fancyping").WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.EndpointPingPOST)).
		WithHandleFuncs("/ping", http.MethodGet, implementation.PingGET).
		WithHandleFuncs("/ping", http.MethodPost, implementation.PingPOST).
		WithGetEndpoint(tamarin.NewEndpoint("/wallet/{}").WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.PrintURLWithElements)).
		WithGetEndpoint(tamarin.NewEndpoint("/wallet/{}/something/{}/else").WithHandlers(implementation.PrintURLWithElements)).
		WithPostEndpoint(tamarin.NewEndpoint("/failjson").WithHandlers(implementation.FailIfNoBody, implementation.ShowBody)).
		WithPostEndpoint(tamarin.NewEndpoint("/content/{*}").WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.StaticSiteHandler)).
		WithGetEndpoint(tamarin.NewEndpoint("/content/{*}").WithHandlers(implementation.StaticSiteHandler))

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
