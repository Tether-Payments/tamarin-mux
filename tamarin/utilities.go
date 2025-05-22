package tamarin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetRequestBodyAndHeader(req *http.Request) ([]byte, http.Header, error) {
	if req == nil {
		return nil, nil, errors.New("nil request")
	}
	if req.Body == nil {
		return nil, nil, errors.New("nil request body")
	}
	if req.Header == nil {
		return nil, nil, errors.New("nil request header")
	}
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read request body : %v", err)
	}
	return bodyBytes, req.Header, nil
}

func UnmarshallJSONRequestBodyTo[T any](req *http.Request, target T) (*T, error) {
	body, _, err := GetRequestBodyAndHeader(req)
	if err != nil {
		return nil, fmt.Errorf("invalid request : %v", err)
	}
	err = json.Unmarshal(body, &target)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal request body to target type : %v", err)
	}
	return &target, nil
}
