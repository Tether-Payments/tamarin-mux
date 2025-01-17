package tamarin

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestNewEndpoint(t *testing.T) {
	ep := NewEndpoint("", "")
	if ep == nil || ep.sequence == nil {
		t.Fail()
	}
}

func TestWithHandlers(t *testing.T) {
	ep := NewEndpoint("", "").WithHandlers(nil)
	if ep == nil {
		t.Fail()
	}
	ep = NewEndpoint("", "").WithHandlers(func(w http.ResponseWriter, r *http.Request) *EndpointError { return nil })
	if len(ep.sequence) != 1 {
		t.Fail()
	}
}

func TestHandle(t *testing.T) {
	testFuncGood := func(rw http.ResponseWriter, req *http.Request) *EndpointError {
		rw.WriteHeader(999)
		rw.Write([]byte("passed"))
		return nil
	}
	testLastCode = -1
	testLastMessage = ""
	ep := NewEndpoint("/test", http.MethodGet).WithHandlers(testFuncGood)
	ep.Handle(testingResponseWriter{}, &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/test"}})
	if testLastCode != 999 || testLastMessage != "passed" {
		t.Fail()
	}
	testFuncBad := func(rw http.ResponseWriter, req *http.Request) *EndpointError {
		return FailWithErrorMessage(-666, "you are a bad person", fmt.Errorf("failing on purpose"))
	}
	testLastCode = -1
	testLastMessage = ""
	ep = NewEndpoint("/test", http.MethodGet).WithHandlers(testFuncBad)
	ep.Handle(testingResponseWriter{}, &http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/test"}})
	if testLastCode != -666 || testLastMessage != "you are a bad person" {
		t.Fail()
	}
}
