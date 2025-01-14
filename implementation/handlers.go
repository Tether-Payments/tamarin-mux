package implementation

import (
	"fmt"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/mux"
)

func PingGET(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (GET)"))
}
func PingPOST(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (POST)"))
}

func EndpointPingPOST(rw http.ResponseWriter, req *http.Request) *mux.EndpointError {
	_, err := rw.Write([]byte("Pong (POST)"))
	if err != nil {
		mux.FailWithErrorMessage(http.StatusInternalServerError, "couldn't write", fmt.Errorf("error while writing response : %v", err))
	}
	return nil
}
