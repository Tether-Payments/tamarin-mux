package tamarin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type EndpointError struct {
	error
	returnCode    int
	returnMessage string
}

func FailWithErrorMessage(code int, message string, err error) *EndpointError {
	return &EndpointError{error: err, returnCode: code, returnMessage: message}
}

func FailWithJSONStatus(code int, v any, err error) *EndpointError {
	jsonBytes, jErr := json.Marshal(v)
	if jErr != nil {
		err = fmt.Errorf("failed to marshall response JSON : %v. Original error was : %v", jErr, err)
		code = http.StatusInternalServerError
	}
	return &EndpointError{error: err, returnCode: code, returnMessage: string(jsonBytes)}
}

func SuceedWithJSONStatus(v any, rw http.ResponseWriter) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Printf("unable to marshal success payload : %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonBytes)
}
