package tamarin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	// staticPaths          []string // TODO : make this a map instead of slice to disconnect actual filesystem from URLS
	handleFuncsGET       map[string][]http.HandlerFunc
	handleFuncsPOST      map[string][]http.HandlerFunc
	variableHandlersGET  map[string][]http.HandlerFunc
	variableHandlersPOST map[string][]http.HandlerFunc
	staticHandlersGET    map[string][]http.HandlerFunc
	staticHandlersPOST   map[string][]http.HandlerFunc
}

func NewHandler() *handler {
	return &handler{
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
		} else if pathIsStaticNew(e.path) {
			s.staticHandlersGET[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsGET[e.path] = []http.HandlerFunc{e.Handle}
		}
	case http.MethodPost:
		if pathIsVariable(e.path) {
			s.variableHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else if pathIsStaticNew(e.path) {
			s.staticHandlersPOST[e.path] = []http.HandlerFunc{e.Handle}
		} else {
			s.handleFuncsPOST[e.path] = []http.HandlerFunc{e.Handle}
		}
	default:
		log.Printf("Don't yet handle the HTTP Method '%s'", e.method)
	}

	return s
}

// func (s *handler) WithStaticDir(path string) *handler {
// 	s.staticPaths = append(s.staticPaths, path)
// 	return s
// }

func (s *handler) WithHandleFuncs(path, httpMethod string, handlerFuncs ...http.HandlerFunc) *handler {
	switch httpMethod {
	case http.MethodGet:
		if pathIsVariable(path) {
			s.variableHandlersGET[path] = handlerFuncs
		} else if pathIsStaticNew(path) {
			s.staticHandlersGET[path] = handlerFuncs
		} else {
			s.handleFuncsGET[path] = handlerFuncs
		}
	case http.MethodPost:
		if pathIsVariable(path) {
			s.variableHandlersPOST[path] = handlerFuncs
		} else if pathIsStaticNew(path) {
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
		// if s.pathIsStatic(reqPath) {
		// 	log.Printf("Using static handler for '%s' ", reqPath)
		// 	endpoints = append(endpoints, s.defaultStaticFunc)
		// 	OK = true
		// } else {
		endpoints, OK = s.handleFuncsGET[reqPath]
		// }
	case http.MethodPost:
		endpoints, OK = s.handleFuncsPOST[reqPath]
	}
	if !OK {
		endpoints = s.getHandlerFuncsForPatternVariable(req.URL.Path, req.Method)
		if endpoints == nil {
			endpoints = s.getHandlerFuncsForPatternStatic(req.URL.Path, req.Method)
			if endpoints == nil {
				log.Printf("don't have a handler for %s", reqPath)
				return
			}
		}
	}
	for _, endpoint := range endpoints {
		endpoint(rw, req)
	}
	log.Printf("Handled request for '%s'", reqPath)
}

// func (h *handler) defaultStaticFunc(rw http.ResponseWriter, req *http.Request) {
// 	log.Printf("[Default Static] Using static handler for '%s'", req.URL.Path)
// 	staticPath, filePath := h.pathComponents(req.URL.Path)
// 	log.Printf("[Default Static] staticpath : %s filepath : %s", staticPath, filePath)
// 	http.ServeFile(rw, req, staticPath+filePath)
// }

// func (h *handler) pathComponents(path string) (string, string) {
// 	staticPortion := ""
// 	minusRoot := ""

// 	for _, staticPath := range h.staticPaths {
// 		trimmed := trimToFirstSlash(staticPath)
// 		if len(path) >= len(trimmed) && path[:len(trimmed)] == trimmed {
// 			staticPortion = staticPath
// 			minusRoot = path[len(trimmed):]
// 		}
// 	}
// 	return staticPortion, minusRoot
// }

// func trimToFirstSlash(path string) string {
// 	idx := strings.Index(path, "/")
// 	if idx < 0 {
// 		return path
// 	}
// 	return path[idx:]
// }

// func (s *handler) pathIsStatic(path string) bool {
// 	for _, candidate := range s.staticPaths {
// 		trimmed := trimToFirstSlash(candidate)
// 		log.Printf("[Path is Static?] Does %s match %s ? ", trimmed, path)
// 		if len(path) >= len(trimmed) && path[:len(trimmed)] == trimmed {
// 			log.Printf("[Path is Static?] %s matches %s", trimmed, path)
// 			return true
// 		}
// 	}
// 	return false
// }

func pathIsStaticNew(path string) bool {
	return strings.Contains(path, "{*}")
}

func pathIsVariable(path string) bool {
	return strings.Contains(path, "{}")
}

func (h *handler) getHandlerFuncsForPatternVariable(path, httpMethod string) []http.HandlerFunc {
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
				if element == "{}" {
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

func (h *handler) getHandlerFuncsForPatternStatic(path, httpMethod string) []http.HandlerFunc {
	log.Printf("[Get HandlerFuncs for Pattern (Static)] Input Path %s", path)
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
			log.Printf("[Get HandlerFuncs for Pattern (Static)] prefix '%s' matches the start of Input Path %s", candidatePrefix, path)
			return handlers

			// candidateSplit := strings.Split(candidatePath, "/")
			// inputSplit := strings.Split(path, "/")
			// if len(candidateSplit) != len(inputSplit) {
			// 	continue
			// }
			// _ = handlers

		}
	}
	return nil
}

func staticPrefix(input string) string {
	idx := strings.Index(input, "{*}")
	if idx < 1 {
		return input
	}
	return input[:idx]
}

func variablePrefix(input string) string {
	idx := strings.Index(input, "{}")
	if idx < 1 {
		return input
	}
	return input[:idx]
}
