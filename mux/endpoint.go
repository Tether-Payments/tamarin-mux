package mux

import (
	"fmt"
	"log"
	"net/http"
)

type EndpointHandlerFunc func(http.ResponseWriter, *http.Request) *EndpointError

type Endpoint struct {
	path     string
	method   string
	sequence []EndpointHandlerFunc
}

func NewEndpoint(path, httpMethod string) *Endpoint {
	return &Endpoint{sequence: []EndpointHandlerFunc{}, path: path, method: httpMethod}
}

func (e *Endpoint) WithHandlers(eFunc ...EndpointHandlerFunc) *Endpoint {
	e.sequence = append(e.sequence, eFunc...)
	return e
}

func (e *Endpoint) Handle(rw http.ResponseWriter, req *http.Request) {
	for _, f := range e.sequence {
		err := f(rw, req)
		if err != nil {
			rw.WriteHeader(err.returnCode)
			rw.Write([]byte(fmt.Sprintf("Bad : %v", err.returnMessage)))
			log.Printf("Exiting due to error %v", err.error)
			break
		}
	}
}
