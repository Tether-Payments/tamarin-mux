package tamarin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	handleFuncsGET       map[string][]http.HandlerFunc
	handleFuncsPOST      map[string][]http.HandlerFunc
	variableHandlersGET  map[string][]http.HandlerFunc
	variableHandlersPOST map[string][]http.HandlerFunc
}

func NewHandler() *handler {
	return &handler{
		handleFuncsGET:       make(map[string][]http.HandlerFunc),
		handleFuncsPOST:      make(map[string][]http.HandlerFunc),
		variableHandlersGET:  make(map[string][]http.HandlerFunc),
		variableHandlersPOST: make(map[string][]http.HandlerFunc),
	}
}

func (s *handler) WithEndpoint(e *Endpoint) *handler {
	switch e.method {
	case http.MethodGet:
		if pathIsVariable(e.path) {
			s.variableHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsGET[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodPost:
		if pathIsVariable(e.path) {
			s.variableHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsPOST[e.path] = []http.HandlerFunc{e.Handle}
		}
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", e.method)
	}

	return s
}

// TODO : Handle Static Pages (http.FileServer?)
func (s *handler) WithStaticDir(path string) *handler {
	return s
}

func (s *handler) WithHandleFuncs(path, httpMethod string, handlerFuncs ...http.HandlerFunc) *handler {
	switch httpMethod {
	case http.MethodGet:
		if pathIsVariable(path) {
			s.variableHandlersGET[path] = handlerFuncs
		} else {
			s.handleFuncsGET[path] = handlerFuncs
		}
	case http.MethodPost:
		if pathIsVariable(path) {
			s.variableHandlersPOST[path] = handlerFuncs
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
	for key := range s.handleFuncsPOST {
		names = append(names, fmt.Sprintf("[%s] -> %s", http.MethodPost, key))
	}
	for key := range s.variableHandlersPOST {
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
		endpoints = s.pathMatchesPattern(req.URL.Path, req.Method)
		if endpoints == nil {
			log.Printf("don't have a handler for %s", reqPath)
			return
		}
	}
	for _, endpoint := range endpoints {
		endpoint(rw, req)
	}
	log.Printf("Handled request for '%s'", reqPath)
}

func pathIsVariable(path string) bool {
	return strings.Contains(path, "{}")
}

func (h *handler) pathMatchesPattern(path, httpMethod string) []http.HandlerFunc {
	var candidateFuncs map[string][]http.HandlerFunc
	switch httpMethod {
	case http.MethodGet:
		log.Printf("potential is GET")
		candidateFuncs = h.variableHandlersGET
	case http.MethodPost:
		log.Printf("potential is POST")
		candidateFuncs = h.variableHandlersPOST
	default:
		return nil
	}
	for candidatePath, handlers := range candidateFuncs {
		rp := variablePrefix(candidatePath)
		// log.Printf("prefix : %s", rp)
		// log.Printf("Path len : %d Prefix len %d", len(path), len(rp))
		if len(path) < len(rp) {
			continue
		}
		// log.Printf("path[:%d] : %s", len(rp), path[:len(rp)])
		if strings.EqualFold(rp, path[:len(rp)]) {
			// log.Println("potential match")
			candidateSplit := strings.Split(candidatePath, "/")
			inputSplit := strings.Split(path, "/")
			// log.Printf("Candidate :%v", candidateSplit)
			// log.Printf("Input     :%v", inputSplit)
			if len(candidateSplit) != len(inputSplit) {
				continue
			}
			allMatched := true
			for idx, element := range candidateSplit {
				if element == "{}" {
					continue
				}
				allMatched = allMatched && strings.EqualFold(element, inputSplit[idx])
			}
			// log.Printf("All matched :%t", allMatched)
			if allMatched {
				return handlers
			}
		}
	}
	return nil
}

func variablePrefix(input string) string {
	idx := strings.Index(input, "{}")
	if idx < 1 {
		return input
	}
	return input[:idx]
}
