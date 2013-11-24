package slingclient

import (
	"log"
	"net/http"
	"os"
	"net/url"
	"path"
	"mime/multipart"
	"bytes"
	"io"
	"path/filepath"
)

type SlingClient struct {

	// Username for authentication
	username string

	// Password for authentication
	password string

	// URL of the sever
	server   url.URL

	// Root path within the server URL space in which this client operates
	rootPath string

	// HTTP client
	httpClient *http.Client

}

// Create a new sling client instance using the given values.
func NewSlingClient(server url.URL, rootPath string, username string, password string) *SlingClient {

	client := &http.Client{}
	return &SlingClient{server: server, username: username, password: password, httpClient: client, rootPath: rootPath}

}

// Creates or updates a node at the given path with the given primary type and properties.
func (slingclient *SlingClient) CreateOrUpdateNode(path string, primaryType string, properties map[string]string) error {

	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)

	writer.WriteField("jcr:primaryType", primaryType)
	for key := range properties {
		log.Print(key)
		writer.WriteField(key, properties[key])
	}

	writer.Close()

	slingclient.PerformMultiPartRequest(path, writer.FormDataContentType(), buf)

	return nil
}


func (slingClient *SlingClient) UploadFile(jcrPath string, file *os.File) error {

	filename := filepath.Base(file.Name())

	// TODO: check filesize; if over threshold use temporary file instead of in-memory buffer
	buf := new(bytes.Buffer)

	writer := multipart.NewWriter(buf)
	// Create form field and get writer to write binary data into:
	w, _ := writer.CreateFormFile(filename, file.Name())

	_, err := io.Copy(w, file)
	if err != nil {
		log.Panic("Error reading file")
	}

	// Write a type hint so sling will create the appropriate node structures
	writer.WriteField(filename + "@TypeHint", "nt:file")

	writer.Close()

	// We need to take the parent path of what we were actually writing to, as Sling will automatically
	// take the name of the file in the multipart request and use that as the last segment of the path.
	effectivePath := filepath.Dir(jcrPath)

	slingClient.PerformMultiPartRequest(effectivePath, writer.FormDataContentType(), buf)
	return nil
}

func (slingclient *SlingClient) DeletePath(jcrPath string) error {

	url, _ := url.ParseRequestURI(path.Join(slingclient.rootPath, jcrPath))
	effectiveUrl := slingclient.server.ResolveReference(url)


	req, err := http.NewRequest("DELETE", effectiveUrl.String(), nil)
	if err != nil {
		log.Panic("Could not construct request: %s", err)
	}
	req.SetBasicAuth(slingclient.username, slingclient.password)

	log.Printf("DELETE %s", jcrPath)

	resp, err := slingclient.httpClient.Do(req)
	if err != nil {
		log.Panicf("Error during request: %s", err)
	}

	log.Printf("       %s", resp.Status)

	return nil
}

func (slingclient *SlingClient) PerformMultiPartRequest(jcrPath string, contentType string, body *bytes.Buffer) error {
	url, _ := url.ParseRequestURI(path.Join(slingclient.rootPath, jcrPath))
	effectiveUrl := slingclient.server.ResolveReference(url)

	log.Printf("POST   %s", jcrPath)

	req, _ := http.NewRequest("POST", effectiveUrl.String(), body)
	req.Header.Add("Content-Type", contentType)
	req.SetBasicAuth(slingclient.username, slingclient.password)

	resp, err := slingclient.httpClient.Do(req)

	if err != nil {
		log.Panicf("Error response")
	}

	log.Printf("       %s", resp.Status)


	return nil
}
