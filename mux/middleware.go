package mux

import (
	"fmt"
	"log"
	"net/http"
)

func MustHaveHelloGoodbyeHeader(rw http.ResponseWriter, req *http.Request) {
	if req.Header == nil {
		failWithError(http.StatusBadRequest, "missing header", rw)
		return
	}
	log.Println("[MiddleWare] Passed first Check : non-nil header")
	val, OK := req.Header["Hello"]
	if !OK || len(val) != 1 {
		failWithError(http.StatusBadRequest, "missing Hello:goodbye in header", rw)
		return
	}
	log.Println("[MiddleWare] Passed second Check : non-nil header")
	if val[0] != "goodbye" {
		failWithError(http.StatusBadRequest, "missing 'goodbye'", rw)
		return
	}
	fmt.Println("[MiddleWare] did not fail")
	rw.WriteHeader(http.StatusOK)
}

func failWithError(code int, message string, rw http.ResponseWriter) {
	rw.WriteHeader(code)
	rw.Write([]byte(message))
}
