package implementation

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Tether-Payments/tamarin-mux/tamarin"
)

func PingGET(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (GET)"))
}

func PingPOST(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("Pong (POST)"))
}

type URLElementsResponse struct {
	Elements map[int]string
}

func PrintURLWithElements(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	split := strings.Split(req.URL.Path, "/")
	elements := make(map[int]string, len(split))
	for i, val := range split {
		elements[i] = val
	}
	tamarin.SuceedWithJSONStatus(URLElementsResponse{Elements: elements}, rw)
	return nil
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

func StaticSiteHandler(rw http.ResponseWriter, req *http.Request) *tamarin.EndpointError {
	log.Printf("[Endpoint StaticSiteHandler] Using static handler for '%s'", req.URL.Path)
	fileSystemPath := mapURLPrefixToFileSystem(req.URL.Path, "/content/", "./static/")
	log.Printf("[Endpoint StaticSiteHandler] FileSystemPath '%s'", fileSystemPath)
	http.ServeFile(rw, req, fileSystemPath)
	return nil
}

func mapURLPrefixToFileSystem(fullPath, urlPrefix, fileSystemPrefix string) string {
	if len(urlPrefix) > len(fullPath) {
		return fullPath
	}
	return fileSystemPrefix + fullPath[len(urlPrefix):]
}
