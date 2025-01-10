package mux

import (
	"net/http"
)

func PingGET(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (GET)"))
}

func PingPOST(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (POST)"))
}
