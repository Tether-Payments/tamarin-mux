package implementation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/tamarin"
)

type ErrorBody struct {
	Message string
}

func FailIfNoBody(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return tamarin.FailWithErrorMessage(http.StatusInternalServerError, "something went wrong on our side", fmt.Errorf("Something went wrong reading request body : %v", err))
	}
	if len(bodyBytes) < 2 {
		return tamarin.FailWithJSONStatus(http.StatusBadRequest, &ErrorBody{Message: "ya done goofed"}, errors.New("no body or body was less than 2 bytes"))
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return nil
}

func MustHaveHelloGoodbyeHeader(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	if req.Header == nil {
		return tamarin.FailWithErrorMessage(http.StatusNoContent, "missing header", fmt.Errorf("missing header"))
	}
	log.Println("[MiddleWare] Passed first Check : non-nil header")

	val, OK := req.Header["Hello"]
	if !OK || len(val) != 1 {
		return tamarin.FailWithErrorMessage(http.StatusBadRequest, "malformed header", fmt.Errorf("missing 'Hello' key in header"))
	}
	log.Println("[MiddleWare] Passed second Check : Hello key in header")

	if val[0] != "goodbye" {
		return tamarin.FailWithErrorMessage(http.StatusBadRequest, "malformed header", fmt.Errorf("missing 'goodbye' value for 'Hello' key"))
	}
	log.Println("[MiddleWare] Passed third Check : 'goodbye' value at 'Hello' key")

	fmt.Println("[MiddleWare] did not fail")
	rw.WriteHeader(http.StatusOK)
	return nil
}
