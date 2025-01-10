package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/mux"
)

var port int

func main() {
	server := mux.NewServer().
		WithEndpoint("/ping", http.MethodGet, mux.PingGET).
		WithEndpoint("/ping", http.MethodPost, mux.PingPOST).
		WithEndpoint("/fancyping", http.MethodPost, mux.NewEndpoint().WithHandlers(mux.MustHaveHelloGoodbyeHeader, mux.PingPOST).HandleFunc)

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
