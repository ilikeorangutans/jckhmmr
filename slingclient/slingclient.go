package slingclient

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"net/url"
	"path"
	"strings"
)

type SlingClient struct {

	// Username for authentication
	Username string

	// Password for authentication
	Password string

	// URL of the sever
	Server   url.URL

	// Root path within the server URL space in which this client operates
	RootPath string

	// HTTP client
	HttpClient *http.Client

	// Local base directory against which all file paths will be relativized.
	BaseDirectory string
}

// Create a new sling client instance using the given values.
func NewSlingClient(server url.URL, rootPath string, username string, password string, baseDirectory string) *SlingClient {

	dir, err := os.Open(baseDirectory)
	if err != nil {
		log.Fatalf("Base directory %s could not be read", baseDirectory)
	}

	dirPath, _ := filepath.Abs(dir.Name())
	log.Printf("Base directory is %s ", dirPath)

	client := &http.Client{}
	return &SlingClient{Server: server, Username: username, Password: password, HttpClient: client, RootPath: rootPath, BaseDirectory: dirPath}

}

func (slingClient *SlingClient) UploadFile(file *os.File) error {

	absoluteFilePath, _ := filepath.Abs(file.Name())
	relPath, err := filepath.Rel(slingClient.BaseDirectory, absoluteFilePath)
	if err != nil {
		log.Fatalf("Cannot relativize %s under %s", file.Name(), slingClient.BaseDirectory)
	}

	segments := strings.Split(relPath, "/")
	for s := range segments {
		log.Print(segments[s])
	}

	url, _ := url.ParseRequestURI(path.Join(slingClient.RootPath, relPath))
	effectiveUrl := slingClient.Server.ResolveReference(url)

	log.Printf("Effective URL %s", effectiveUrl.String())
	req, _ := http.NewRequest("GET", effectiveUrl.String(), nil)

	resp, err := slingClient.HttpClient.Do(req)
	if err != nil {
		log.Fatalf("Error during GET %s", err)
	}

	log.Print(resp)

	return nil
}
