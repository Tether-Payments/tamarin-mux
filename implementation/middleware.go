package implementation

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tether-Payments/tamarin-mux/mux"
)

func MustHaveHelloGoodbyeHeader(rw http.ResponseWriter, req *http.Request) *mux.EndpointError {
	if req.Header == nil {
		return mux.FailWithErrorMessage(http.StatusNoContent, "missing header", fmt.Errorf("missing header"))
	}
	log.Println("[MiddleWare] Passed first Check : non-nil header")
	val, OK := req.Header["Hello"]
	if !OK || len(val) != 1 {
		return mux.FailWithErrorMessage(http.StatusBadRequest, "malformed header", fmt.Errorf("missing 'Hello' key in header"))
	}
	log.Println("[MiddleWare] Passed second Check : Hello key in header")
	if val[0] != "goodbye" {
		return mux.FailWithErrorMessage(http.StatusBadRequest, "malformed header", fmt.Errorf("missing 'goodbye' value for 'Hello' key"))
	}
	log.Println("[MiddleWare] Passed third Check : 'goodbye' value at 'Hello' key")
	fmt.Println("[MiddleWare] did not fail")
	rw.WriteHeader(http.StatusOK)
	return nil
}
