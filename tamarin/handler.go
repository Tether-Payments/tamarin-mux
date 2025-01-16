package tamarin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	VARIABLE_INDICATOR = "{}"
	STATIC_INDICATOR   = "{*}"
)

type handler struct {
	verbose              bool
	handleFuncsGET       map[string][]http.HandlerFunc
	handleFuncsPOST      map[string][]http.HandlerFunc
	variableHandlersGET  map[string][]http.HandlerFunc
	variableHandlersPOST map[string][]http.HandlerFunc
	staticHandlersGET    map[string][]http.HandlerFunc
	staticHandlersPOST   map[string][]http.HandlerFunc
}

func NewHandler(verbose bool) *handler {
	return &handler{
		verbose:              verbose,
		handleFuncsGET:       make(map[string][]http.HandlerFunc),
		handleFuncsPOST:      make(map[string][]http.HandlerFunc),
		variableHandlersGET:  make(map[string][]http.HandlerFunc),
		variableHandlersPOST: make(map[string][]http.HandlerFunc),
		staticHandlersGET:    make(map[string][]http.HandlerFunc),
		staticHandlersPOST:   make(map[string][]http.HandlerFunc),
	}
}

func (s *handler) WithEndpoint(e *Endpoint) *handler {
	switch e.method {
	case http.MethodGet:
		if pathIsVariable(e.path) {
			s.variableHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsGET[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodPost:
		if pathIsVariable(e.path) {
			s.variableHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStatic(e.path) {
			s.staticHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsPOST[e.path] = []http.HandlerFunc{e.Handle}
		}
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", e.method)
	}

	return s
}

func (s *handler) WithHandleFuncs(path, httpMethod string, handlerFuncs ...http.HandlerFunc) *handler {
	switch httpMethod {
	case http.MethodGet:
		if pathIsVariable(path) {
			s.variableHandlersGET[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersGET[path] = handlerFuncs
		} else {
			s.handleFuncsGET[path] = handlerFuncs
		}
	case http.MethodPost:
		if pathIsVariable(path) {
			s.variableHandlersPOST[path] = handlerFuncs
		} else if pathIsStatic(path) {
			s.staticHandlersPOST[path] = handlerFuncs
		} else {
			s.handleFuncsPOST[path] = handlerFuncs
		}
		s.handleFuncsPOST[path] = handlerFuncs
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", httpMethod)
	}
	return s
}

func (s *handler) HandlerNames() []string {
	names := []string{}
	for key := range s.handleFuncsGET {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodGet, key))
	}
	for key := range s.variableHandlersGET {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodGet, key))
	}
	for key := range s.staticHandlersGET {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodGet, key))
	}
	for key := range s.handleFuncsPOST {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodPost, key))
	}
	for key := range s.variableHandlersPOST {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodPost, key))
	}
	for key := range s.staticHandlersPOST {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodPost, key))
	}
	return names
}

func (s *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	log.Printf("Received request for '%s'", reqPath)
	var endpoints []http.HandlerFunc
	var OK bool
	switch req.Method {
	case http.MethodGet:
		endpoints, OK = s.handleFuncsGET[reqPath]
	case http.MethodPost:
		endpoints, OK = s.handleFuncsPOST[reqPath]
	}
	if !OK {
		endpoints = s.getVariableHandlerFuncsForPattern(req.URL.Path, req.Method)
		if endpoints == nil {
			endpoints = s.getStaticHandlerFuncsForPattern(req.URL.Path, req.Method)
			if endpoints == nil {
				if s.verbose {
					log.Printf("don't have a handler for %s", reqPath)
				}
				return
			}
		}
	}
	for _, endpoint := range endpoints {
		endpoint(rw, req)
	}
	if s.verbose {
		log.Printf("Handled request for '%s'", reqPath)
	}
}

func pathIsStatic(path string) bool {
	return strings.Contains(path, STATIC_INDICATOR)
}

func pathIsVariable(path string) bool {
	return strings.Contains(path, VARIABLE_INDICATOR)
}

func (h *handler) getVariableHandlerFuncsForPattern(path, httpMethod string) []http.HandlerFunc {
	var candidateFuncs map[string][]http.HandlerFunc
	switch httpMethod {
	case http.MethodGet:
		candidateFuncs = h.variableHandlersGET
	case http.MethodPost:
		candidateFuncs = h.variableHandlersPOST
	default:
		return nil
	}
	for candidatePath, handlers := range candidateFuncs {
		candidatePrefix := variablePrefix(candidatePath)
		if len(path) < len(candidatePrefix) {
			continue
		}
		if strings.EqualFold(candidatePrefix, path[:len(candidatePrefix)]) {
			candidateSplit := strings.Split(candidatePath, "/")
			inputSplit := strings.Split(path, "/")
			if len(candidateSplit) != len(inputSplit) {
				continue
			}
			allMatched := true
			for idx, element := range candidateSplit {
				if element == VARIABLE_INDICATOR {
					continue
				}
				allMatched = allMatched && strings.EqualFold(element, inputSplit[idx])
			}
			if allMatched {
				return handlers
			}
		}
	}
	return nil
}

func (h *handler) getStaticHandlerFuncsForPattern(path, httpMethod string) []http.HandlerFunc {
	var candidateFuncs map[string][]http.HandlerFunc
	switch httpMethod {
	case http.MethodGet:
		candidateFuncs = h.staticHandlersGET
	case http.MethodPost:
		candidateFuncs = h.staticHandlersPOST
	default:
		return nil
	}
	fmt.Println(len(candidateFuncs))
	for candidatePath, handlers := range candidateFuncs {
		candidatePrefix := staticPrefix(candidatePath)
		if len(path) < len(candidatePrefix) {
			continue
		}
		if strings.EqualFold(candidatePrefix, path[:len(candidatePrefix)]) {
			return handlers
		}
	}
	return nil
}

func staticPrefix(input string) string {
	idx := strings.Index(input, STATIC_INDICATOR)
	if idx < 1 {
		return input
	}
	return input[:idx]
}

func variablePrefix(input string) string {
	idx := strings.Index(input, VARIABLE_INDICATOR)
	if idx < 1 {
		return input
	}
	return input[:idx]
}
