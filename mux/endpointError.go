package mux

type EndpointError struct {
	error
	returnCode    int
	returnMessage string
}

func FailWithErrorMessage(code int, message string, err error) *EndpointError {
	return &EndpointError{error: err, returnCode: code, returnMessage: message}
}
