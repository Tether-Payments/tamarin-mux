package tamarin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(false)
	if h == nil {
		t.Fail()
	}
}

func TestWithEndpoint(t *testing.T) {
	h := NewHandler(false).withEndpoint(nil)
	if h == nil {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{})
	if h == nil {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodGet})
	if h == nil {
		t.Fail()
	}
	testFunc := func(http.ResponseWriter, *http.Request) *EndpointError { return nil }
	h = h.withEndpoint(&endpoint{method: http.MethodGet, sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.handleFuncsGET) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodGet, path: "/{}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.variableHandlersGET) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodGet, path: "/{*}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.staticHandlersGET) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPost, sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.handleFuncsPOST) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPost, path: "/{}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.variableHandlersPOST) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPost, path: "/{*}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.staticHandlersPOST) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPatch, sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.handleFuncsPATCH) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPatch, path: "/{}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.variableHandlersPATCH) != 1 {
		t.Fail()
	}
	h = h.withEndpoint(&endpoint{method: http.MethodPatch, path: "/{*}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.staticHandlersPATCH) != 1 {
		t.Fail()
	}
}

func TestWithGetEndpoint(t *testing.T) {
	testFunc := func(http.ResponseWriter, *http.Request) *EndpointError { return nil }
	h := NewHandler(true).WithGetEndpoint(NewEndpoint("").WithHandlers(testFunc))
	if len(h.handleFuncsGET) != 1 {
		t.Fail()
	}
	h = NewHandler(true).WithGetEndpoint(nil)
	if len(h.handleFuncsGET) != 0 {
		t.Fail()
	}
}

func TestWithPostEndpoint(t *testing.T) {
	testFunc := func(http.ResponseWriter, *http.Request) *EndpointError { return nil }
	h := NewHandler(true).WithPostEndpoint(NewEndpoint("").WithHandlers(testFunc))
	if len(h.handleFuncsPOST) != 1 {
		t.Fail()
	}
	h = NewHandler(true).WithPostEndpoint(nil)
	if len(h.handleFuncsPOST) != 0 {
		t.Fail()
	}
}

func TestWithPatchEndpoint(t *testing.T) {
	testFunc := func(http.ResponseWriter, *http.Request) *EndpointError { return nil }
	h := NewHandler(true).WithPatchEndpoint(NewEndpoint("").WithHandlers(testFunc))
	if len(h.handleFuncsPATCH) != 1 {
		t.Fail()
	}
	h = NewHandler(true).WithPatchEndpoint(nil)
	if len(h.handleFuncsPATCH) != 0 {
		t.Fail()
	}
}

func TestWithHandleFuncs(t *testing.T) {
	h := NewHandler(false).WithHandleFuncs("", "", nil)
	if h == nil {
		t.Fail()
	}
	h = NewHandler(false).WithHandleFuncs("", http.MethodGet, nil)
	if h == nil {
		t.Fail()
	}
	testFunc := func(http.ResponseWriter, *http.Request) {}
	h = h.WithHandleFuncs("", http.MethodGet, testFunc)
	if len(h.handleFuncsGET) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("/{}", http.MethodGet, testFunc)
	if len(h.variableHandlersGET) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("/{*}", http.MethodGet, testFunc)
	if len(h.staticHandlersGET) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("", http.MethodPost, testFunc)
	if len(h.handleFuncsPOST) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("{}", http.MethodPost, testFunc)
	if len(h.variableHandlersPOST) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("{*}", http.MethodPost, testFunc)
	if len(h.staticHandlersPOST) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("", http.MethodPatch, testFunc)
	if len(h.handleFuncsPATCH) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("{}", http.MethodPatch, testFunc)
	if len(h.variableHandlersPATCH) != 1 {
		t.Fail()
	}
	h = h.WithHandleFuncs("{*}", http.MethodPatch, testFunc)
	if len(h.staticHandlersPATCH) != 1 {
		t.Fail()
	}
}

func TestHandlerNames(t *testing.T) {
	testFunc := func(http.ResponseWriter, *http.Request) {}
	h := NewHandler(false).
		WithHandleFuncs("/pathA", http.MethodGet, testFunc).
		WithHandleFuncs("/pathB/{}", http.MethodGet, testFunc).
		WithHandleFuncs("/pathC/{*}", http.MethodGet, testFunc).
		WithHandleFuncs("/pathD/", http.MethodPost, testFunc).
		WithHandleFuncs("/pathE/{}", http.MethodPost, testFunc).
		WithHandleFuncs("/pathF/{*}", http.MethodPost, testFunc).
		WithHandleFuncs("/pathG/", http.MethodPatch, testFunc).
		WithHandleFuncs("/pathH/{}", http.MethodPatch, testFunc).
		WithHandleFuncs("/pathI/{*}", http.MethodPatch, testFunc)
	names := h.HandlerNames()
	if len(names) != 9 {
		t.Fail()
	}
	if !strings.Contains(names[0], "/pathA") || !strings.Contains(names[0], http.MethodGet) {
		t.Fail()
	}
	if !strings.Contains(names[1], "/pathB") || !strings.Contains(names[1], http.MethodGet) {
		t.Fail()
	}
	if !strings.Contains(names[2], "/pathC") || !strings.Contains(names[2], http.MethodGet) {
		t.Fail()
	}
	if !strings.Contains(names[3], "/pathD") || !strings.Contains(names[3], http.MethodPost) {
		t.Fail()
	}
	if !strings.Contains(names[4], "/pathE") || !strings.Contains(names[4], http.MethodPost) {
		t.Fail()
	}
	if !strings.Contains(names[5], "/pathF") || !strings.Contains(names[5], http.MethodPost) {
		t.Fail()
	}
	if !strings.Contains(names[6], "/pathG") || !strings.Contains(names[6], http.MethodPatch) {
		t.Fail()
	}
	if !strings.Contains(names[7], "/pathH") || !strings.Contains(names[7], http.MethodPatch) {
		t.Fail()
	}
	if !strings.Contains(names[8], "/pathI") || !strings.Contains(names[8], http.MethodPatch) {
		t.Fail()
	}
}

func TestServeHTTP(t *testing.T) {
	testLastMessage = ""
	testLastCode = 0
	h := NewHandler(true)
	h.ServeHTTP(nil, nil)
	trw := testingResponseWriter{}
	h.ServeHTTP(trw, &http.Request{})
	if testLastMessage != "" && testLastCode != 0 {
		t.Fail()
	}
	h.ServeHTTP(trw, &http.Request{URL: &url.URL{Path: "/test"}})
	if testLastCode != http.StatusNotFound {
		t.Fail()
	}
	testLastCode = -1
	h.ServeHTTP(trw, &http.Request{URL: &url.URL{Path: "/test"}, Method: http.MethodGet})
	if testLastCode != http.StatusNotFound {
		t.Fail()
	}
	testLastCode = -1
	h.ServeHTTP(trw, &http.Request{URL: &url.URL{Path: "/test"}, Method: http.MethodPost})
	if testLastCode != http.StatusNotFound {
		t.Fail()
	}
	goodMessage := "you have passed the test"
	goodStatus := 999
	testFunc := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(goodStatus)
		rw.Write([]byte(goodMessage))
	}
	h.WithHandleFuncs("/test", http.MethodGet, testFunc)
	testLastCode = -1
	testLastMessage = ""
	h.ServeHTTP(trw, &http.Request{URL: &url.URL{Path: "/test"}, Method: http.MethodGet})
	if testLastCode != goodStatus {
		t.Fail()
	}
	if testLastMessage != goodMessage {
		t.Fail()
	}
	fmt.Println(testLastMessage, testLastCode)
	h.WithHandleFuncs("/test", http.MethodPatch, testFunc)
	testLastCode = -1
	testLastMessage = ""
	h.ServeHTTP(trw, &http.Request{URL: &url.URL{Path: "/test"}, Method: http.MethodPatch})
	if testLastCode != goodStatus {
		t.Fail()
	}
	if testLastMessage != goodMessage {
		t.Fail()
	}
	fmt.Println(testLastMessage, testLastCode)
}

func TestGetVariableHandlerFuncsForPattern(t *testing.T) {
	h := NewHandler(true)
	result := h.getVariableHandlerFuncsForPattern("", "")
	if result != nil {
		t.Fail()
	}
	goodMessage := "you have passed the test"
	goodStatus := 999
	testFunc := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(goodMessage))
		rw.WriteHeader(goodStatus)
	}
	testLastMessage = ""
	testLastCode = -1
	result = NewHandler(true).WithHandleFuncs("/test/{}", http.MethodGet, testFunc).getVariableHandlerFuncsForPattern("/", http.MethodGet)
	if result != nil {
		t.Fail()
	}
	h = NewHandler(true).WithHandleFuncs("/test/{}/something/else/{}/", http.MethodGet, testFunc).WithHandleFuncs("/test/{}", http.MethodGet, testFunc)
	h.ServeHTTP(testingResponseWriter{}, &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/test/anotherword"}})
	if testLastMessage != goodMessage || testLastCode != goodStatus {
		t.Fail()
	}
	h = NewHandler(true).WithHandleFuncs("/test/{}/something/else/{}/", http.MethodGet, testFunc).WithHandleFuncs("/test/{}", http.MethodPatch, testFunc)
	h.ServeHTTP(testingResponseWriter{}, &http.Request{Method: http.MethodPatch, URL: &url.URL{Path: "/test/anotherword"}})
	if testLastMessage != goodMessage || testLastCode != goodStatus {
		t.Fail()
	}

}

func TestGetStaticVariableHandlerFuncsForPattern(t *testing.T) {
	h := NewHandler(true)
	result := h.getStaticHandlerFuncsForPattern("", "")
	if result != nil {
		t.Fail()
	}
	goodMessage := "you have passed the test"
	goodStatus := 999
	testFunc := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(goodMessage))
		rw.WriteHeader(goodStatus)
	}
	testLastMessage = ""
	testLastCode = -1
	result = NewHandler(true).WithHandleFuncs("/test/{*}", http.MethodGet, testFunc).getVariableHandlerFuncsForPattern("/", http.MethodGet)
	if result != nil {
		t.Fail()
	}
	h = NewHandler(true).WithHandleFuncs("/test/something/else/{*}/", http.MethodGet, testFunc).WithHandleFuncs("/test/{*}", http.MethodGet, testFunc)
	h.ServeHTTP(testingResponseWriter{}, &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/test/myfile.json"}})
	if testLastMessage != goodMessage || testLastCode != goodStatus {
		t.Fail()
	}
	h = NewHandler(true).WithHandleFuncs("/test/something/else/{*}/", http.MethodPatch, testFunc).WithHandleFuncs("/test/{*}", http.MethodPatch, testFunc)
	h.ServeHTTP(testingResponseWriter{}, &http.Request{Method: http.MethodPatch, URL: &url.URL{Path: "/test/myfile.json"}})
	if testLastMessage != goodMessage || testLastCode != goodStatus {
		t.Fail()
	}
}

func TestStaticPrefix(t *testing.T) {
	result := staticPrefix("")
	if len(result) != 0 {
		t.Fail()
	}
	result = staticPrefix("noprefix")
	if len(result) != len("noprefix") {
		t.Fail()
	}
	result = staticPrefix(fmt.Sprintf("word%sanotherword", STATIC_INDICATOR))
	if result != "word" {
		t.Fail()
	}
}
func TestVariablePrefix(t *testing.T) {
	result := variablePrefix("")
	if len(result) != 0 {
		t.Fail()
	}
	result = variablePrefix("noprefix")
	if len(result) != len("noprefix") {
		t.Fail()
	}
	result = variablePrefix(fmt.Sprintf("word%sanotherword", VARIABLE_INDICATOR))
	if result != "word" {
		t.Fail()
	}
}

func TestPost(t *testing.T) {
	m := NewHandler(false)
	m.Post("/test1", func(w http.ResponseWriter, r *http.Request) *EndpointError { return nil })
	if len(m.handleFuncsPOST) != 1 {
		t.Fail()
	}
}

func TestPostF(t *testing.T) {
	m := NewHandler(false).
		PostF("/test1", func(w http.ResponseWriter, r *http.Request) {}).
		PostF("/test2/{}", func(w http.ResponseWriter, r *http.Request) {}).
		PostF("/test3/{*}", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.handleFuncsPOST) != 1 || len(m.variableHandlersPOST) != 1 || len(m.staticHandlersPOST) != 1 {
		t.Fail()
	}
}
func TestGet(t *testing.T) {
	m := NewHandler(false)
	m.Get("/test1", func(w http.ResponseWriter, r *http.Request) *EndpointError { return nil })
	if len(m.handleFuncsGET) != 1 {
		t.Fail()
	}
}
func TestGetF(t *testing.T) {
	m := NewHandler(false)
	m.GetF("/test1", func(w http.ResponseWriter, r *http.Request) {})
	m.GetF("/test2/{}", func(w http.ResponseWriter, r *http.Request) {})
	m.GetF("/test3/{*}", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.handleFuncsGET) != 1 || len(m.variableHandlersGET) != 1 || len(m.staticHandlersGET) != 1 {
		t.Fail()
	}
}
func TestPatch(t *testing.T) {
	m := NewHandler(false)
	m.Patch("/test1", func(w http.ResponseWriter, r *http.Request) *EndpointError { return nil })
	if len(m.handleFuncsPATCH) != 1 {
		t.Fail()
	}
}
func TestPatchF(t *testing.T) {
	m := NewHandler(false)
	m.PatchF("/test1", func(w http.ResponseWriter, r *http.Request) {})
	m.PatchF("/test2/{}", func(w http.ResponseWriter, r *http.Request) {})
	m.PatchF("/test3/{*}", func(w http.ResponseWriter, r *http.Request) {})
	if len(m.handleFuncsPATCH) != 1 || len(m.variableHandlersPATCH) != 1 || len(m.staticHandlersPATCH) != 1 {
		t.Fail()
	}
}

var (
	testLastMessage string
	testLastCode    int
)

type testingResponseWriter struct{}

func (t testingResponseWriter) Header() http.Header {
	return http.Header{}
}

func (t testingResponseWriter) Write(message []byte) (int, error) {
	testLastMessage = string(message)
	return len(message), nil
}

func (t testingResponseWriter) WriteHeader(code int) {
	testLastCode = code
}
