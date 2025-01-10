package mux

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	endpointsGET  map[string][]http.HandlerFunc
	endpointsPOST map[string][]http.HandlerFunc
}

func NewServer() *server {
	return &server{
		endpointsGET:  make(map[string][]http.HandlerFunc),
		endpointsPOST: make(map[string][]http.HandlerFunc),
	}
}

func (s *server) WithEndpoint(path, httpMethod string, handlerFuncs ...http.HandlerFunc) *server {
	switch httpMethod {
	case http.MethodGet:
		s.endpointsGET[path] = handlerFuncs
	case http.MethodPost:
		s.endpointsPOST[path] = handlerFuncs
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", httpMethod)
	}
	return s
}

func (s *server) HandlerNames() []string {
	names := []string{}
	for key := range s.endpointsGET {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodGet, key))
	}
	for key := range s.endpointsPOST {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodPost, key))
	}
	return names
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	log.Printf("Received request for '%s'", reqPath)
	var endpoints []http.HandlerFunc
	var OK bool
	switch req.Method {
	case http.MethodGet:
		endpoints, OK = s.endpointsGET[reqPath]
	case http.MethodPost:
		endpoints, OK = s.endpointsPOST[reqPath]
	}
	if !OK {
		log.Printf("don't have a handler for %s", reqPath)
		return
	}
	for _, endpoint := range endpoints {
		endpoint(rw, req)
	}
	log.Printf("Handled request for '%s'", reqPath)
}
