package tamarin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestGetRequestBodyAndHeader(t *testing.T) {
	_, _, err := GetRequestBodyAndHeader(nil)
	if err == nil {
		t.Fail()
	}
	_, _, err = GetRequestBodyAndHeader(&http.Request{})
	if err == nil {
		t.Fail()
	}
	_, _, err = GetRequestBodyAndHeader(&http.Request{Body: io.NopCloser(bytes.NewReader([]byte{}))})
	if err == nil {
		t.Fail()
	}
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader([]byte("A")))
	if err != nil {
		t.Fail()
	}
	req.Header.Add("B", "C")
	body, header, err := GetRequestBodyAndHeader(req)
	if err != nil {
		t.Fail()
	}
	if !bytes.Equal(body, []byte("A")) {
		t.Fail()
	}
	if header.Get("B") != "C" {
		t.Fail()
	}
}

type testingStructA struct {
	S string
	F float64
}
type testingStructB struct {
	B bool
	I int
}

func TestUnmarshallJSONRequestBodyTo(t *testing.T) {
	tsA := testingStructA{S: "A", F: 1.23}
	bodyBytesA, err := json.Marshal(tsA)
	if err != nil {
		t.Fail()
	}
	req, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(bodyBytesA))
	if err != nil {
		t.Log("err request")
		t.Fail()
	}
	req.Header.Add("B", "C")

	result, _, err := UnmarshallJSONRequestBodyTo(req, testingStructA{})
	if err != nil {
		t.Log("err happy")
		t.Fail()
	}
	if result == nil || result.S != "A" {
		t.Logf("err result : %v", err)
		t.Fail()
	}
	_, _, err = UnmarshallJSONRequestBodyTo(req, "")
	if err == nil {
		t.Logf("err sad : %v", err)
		t.Fail()
	}
}
