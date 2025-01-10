package mux

import (
	"fmt"
	"log"
	"net/http"
)

type server struct {
	endpointsGET  map[string]func(http.ResponseWriter, *http.Request)
	endpointsPOST map[string]func(http.ResponseWriter, *http.Request)
}

func NewServer() *server {
	return &server{
		endpointsGET:  make(map[string]func(http.ResponseWriter, *http.Request)),
		endpointsPOST: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (s *server) WithEndpoint(path, httpMethod string, handlerFunc func(http.ResponseWriter, *http.Request)) *server {
	switch httpMethod {
	case http.MethodGet:
		s.endpointsGET[path] = handlerFunc
	case http.MethodPost:
		s.endpointsPOST[path] = handlerFunc
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
	var endpoint func(http.ResponseWriter, *http.Request)
	var OK bool
	switch req.Method {
	case http.MethodGet:
		endpoint, OK = s.endpointsGET[reqPath]
	case http.MethodPost:
		endpoint, OK = s.endpointsPOST[reqPath]
	}
	if !OK {
		log.Printf("don't have a handler for %s", reqPath)
		return
	}

	endpoint(rw, req)
	log.Printf("Handled request for '%s'", reqPath)
}
