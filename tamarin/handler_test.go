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
	h := NewHandler(false).WithEndpoint(nil)
	if h == nil {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{})
	if h == nil {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodGet})
	if h == nil {
		t.Fail()
	}
	testFunc := func(http.ResponseWriter, *http.Request) *EndpointError { return nil }
	h = h.WithEndpoint(&endpoint{method: http.MethodGet, sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.handleFuncsGET) != 1 {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodGet, path: "/{}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.variableHandlersGET) != 1 {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodGet, path: "/{*}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.staticHandlersGET) != 1 {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodPost, sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.handleFuncsPOST) != 1 {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodPost, path: "/{}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.variableHandlersPOST) != 1 {
		t.Fail()
	}
	h = h.WithEndpoint(&endpoint{method: http.MethodPost, path: "/{*}", sequence: []EndpointHandlerFunc{testFunc}})
	if len(h.staticHandlersPOST) != 1 {
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
}

func TestHandlerNames(t *testing.T) {
	testFunc := func(http.ResponseWriter, *http.Request) {}
	h := NewHandler(false).
		WithHandleFuncs("/pathA", http.MethodGet, testFunc).
		WithHandleFuncs("/pathB/{}", http.MethodGet, testFunc).
		WithHandleFuncs("/pathC/{*}", http.MethodGet, testFunc).
		WithHandleFuncs("/pathD/", http.MethodPost, testFunc).
		WithHandleFuncs("/pathE/{}", http.MethodPost, testFunc).
		WithHandleFuncs("/pathF/{*}", http.MethodPost, testFunc)
	names := h.HandlerNames()
	if len(names) != 6 {
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
