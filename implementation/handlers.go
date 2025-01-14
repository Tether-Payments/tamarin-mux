package implementation

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/tamarin"
)

func PingGET(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (GET)"))
}

func PingPOST(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (POST)"))
}

func EndpointPingPOST(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	_, err := rw.Write([]byte("Pong (POST)"))
	if err != nil {
		tamarin.FailWithErrorMessage(http.StatusInternalServerError, "couldn't write", fmt.Errorf("error while writing response : %v", err))
	}
	return nil
}

type showBodyResponse struct {
	YourBody string
}

func ShowBody(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		tamarin.FailWithErrorMessage(http.StatusInternalServerError, "failed reading body", err)
	}
	log.Printf("Request Body :\n%s", bodyBytes)
	tamarin.SuceedWithJSONStatus(showBodyResponse{YourBody: string(bodyBytes)}, rw)
	return nil
}
