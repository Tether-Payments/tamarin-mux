package tamarin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

type testingResponseBody struct {
	I int
	S string
}

func TestFailWithErrorMessage(t *testing.T) {
	r := FailWithErrorMessage(-1, "", nil)
	if r == nil {
		t.Fail()
	}
	r = FailWithErrorMessage(1, "two", errors.ErrUnsupported)
	if !errors.Is(r.error, errors.ErrUnsupported) || r.returnCode != 1 || r.returnMessage != "two" {
		t.Fail()
	}
}

func TestFailWithJSONStatus(t *testing.T) {
	r := FailWithJSONStatus(-1, nil, nil)
	if r == nil {
		t.Fail()
	}
	inBody := testingResponseBody{I: 100, S: "one hundred"}
	result := FailWithJSONStatus(http.StatusNotFound, inBody, errors.ErrUnsupported)
	if !errors.Is(result.error, errors.ErrUnsupported) || result.returnCode != http.StatusNotFound {
		t.Fail()
	}
	var outBody testingResponseBody
	err := json.Unmarshal([]byte(result.returnMessage), &outBody)
	if err != nil {
		t.Fail()
	}
	if outBody.I != inBody.I || outBody.S != inBody.S {
		t.Fail()
	}
	result = FailWithJSONStatus(-1, http.Request{Body: io.NopCloser(bytes.NewReader([]byte("this should fail")))}, errors.New("original"))
	if result.returnCode != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestSucceedWithJSONStatus(t *testing.T) {
	r := SuceedWithJSONStatus(nil, nil)
	if r.returnCode != http.StatusInternalServerError {
		t.Fail()
	}
	inbody := testingResponseBody{I: 123, S: "one two thre"}
	jb, err := json.Marshal(inbody)
	if err != nil {
		t.Fail()
	}
	testLastCode = -1
	testLastMessage = ""
	r = SuceedWithJSONStatus(inbody, testingResponseWriter{})
	if testLastCode != http.StatusOK {
		t.Fail()
	}
	if testLastMessage != string(jb) {
		t.Fail()
	}
	result := SuceedWithJSONStatus(http.Request{Body: io.NopCloser(bytes.NewReader([]byte("this should fail")))}, testingResponseWriter{})
	if result == nil {
		t.Fail()
	}

}
