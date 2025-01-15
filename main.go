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

const uuidRegex = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`

func main() {
	handler := tamarin.NewHandler().
		WithEndpoint(tamarin.NewEndpoint("/fancyping", http.MethodPost).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.EndpointPingPOST)).
		WithHandleFuncs("/ping", http.MethodGet, implementation.PingGET).
		WithHandleFuncs("/ping", http.MethodPost, implementation.PingPOST).
		WithEndpoint(tamarin.NewEndpoint("/wallet/{}", http.MethodGet).WithHandlers(implementation.MustHaveHelloGoodbyeHeader, implementation.PrintURLWithElements)).
		WithEndpoint(tamarin.NewEndpoint("/wallet/{}/something/{}/else", http.MethodGet).WithHandlers(implementation.PrintURLWithElements)).
		WithEndpoint(tamarin.NewEndpoint("/failjson", http.MethodPost).WithHandlers(implementation.FailIfNoBody, implementation.ShowBody)).
		WithStaticDir("./static")

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
