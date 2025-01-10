package mux

import "net/http"

type Endpoint struct {
	sequence       []http.HandlerFunc
	success        bool
	responseWriter responseWriterWrapper
}

func NewEndpoint() *Endpoint {
	return &Endpoint{sequence: []http.HandlerFunc{}}
}

func (e *Endpoint) WithHandlers(handlers ...http.HandlerFunc) *Endpoint {
	for _, handler := range handlers {
		e.sequence = append(e.sequence, handler)
	}
	return e
}

func (e *Endpoint) HandleFunc(rw http.ResponseWriter, req *http.Request) {
	e.responseWriter = responseWriterWrapper{actualWriter: rw, owner: e}
	for _, next := range e.sequence {
		next(e.responseWriter, req)
		if !e.success {
			return
		}
	}
}

type responseWriterWrapper struct {
	owner        *Endpoint
	actualWriter http.ResponseWriter
}

func (c responseWriterWrapper) Header() http.Header {
	return c.actualWriter.Header()
}

func (c responseWriterWrapper) Write(input []byte) (int, error) {
	return c.actualWriter.Write(input)
}

// TODO : make owning endpoint stop the middleware immediately (to get rid of all the empty "returns")
func (c responseWriterWrapper) WriteHeader(statusCode int) {
	c.owner.success = statusCode == http.StatusOK
	if statusCode != http.StatusOK {
		c.actualWriter.WriteHeader(statusCode)
	}
}
