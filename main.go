package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"net/url"
	"jckhmmr/slingclient"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "jckhmmr"
	app.Version = "0.0.1"
	app.Usage = "Jackhammer written in GO"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "host",
			Value: "http://server:port",
			Usage: "Server to connect to",
		},
		cli.StringFlag{
			Name: "username",
			Value: "admin",
			Usage: "Username",
		},
		cli.StringFlag{
			Name: "password",
			Value: "admin",
			Usage: "Passowrd",
		},
	}

	app.Commands = []cli.Command{

		{
			Name: "upload",
			Usage: "Upload a single file",
			Action: UploadFile,
			Flags: []cli.Flag {

				cli.StringFlag {
					Name: "f",
					Value: "FILE",
					Usage: "File to upload",
				},
				cli.StringFlag {
					Name: "",
					Usage: "asdf",
				},
			},
		},
		{
			Name: "watch",
			Usage: "Watch a directory and automatically mirror all changes",
			Action: WatchDirectory,
		},
	}

	app.Run(os.Args)
}

func WatchDirectory(c *cli.Context) {

	dirPath := c.String("dir")
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Panic("Could not open directory")
	}

	log.Print(dir)
}

func UploadFile(c * cli.Context) {

	//slingclient.NewSlingClient("http://localhost:8080", "/", "admin", "admin", )

	//log.Print("Uploading file")
}


func foo() {


	port := os.Args[1]
	node := os.Args[2]
	fileToUpload := os.Args[3]

	client := &http.Client{}

	serverUrl, _ := url.Parse("http://localhost:8080/")
	slingClient := slingclient.NewSlingClient(*serverUrl, "/tmp", "admin", "admin", "./")


	buf := new(bytes.Buffer)

	nodeName := filepath.Base(fileToUpload)
	log.Printf("final node name %s", nodeName)

	multiPartWriter := multipart.NewWriter(buf)
	w, _ := multiPartWriter.CreateFormFile(fileToUpload, fileToUpload)

	fd, err := os.Open(fileToUpload)
	if err != nil {
		log.Fatalf("Could not open file %s", fileToUpload)
	}
	defer fd.Close()

	slingClient.UploadFile(fd)

	_, err = io.Copy(w, fd)
	if err != nil {
		log.Fatal("Error while copying file data")
		return
	}
	multiPartWriter.WriteField(nodeName + "@TypeHint", "nt:file")

	multiPartWriter.Close()
	//log.Print(buf)

	req, err := http.NewRequest("POST", "http://localhost:" + port + "/" + node, buf)
	if err != nil {
		log.Panic("Error creating requerst")
	}
	req.SetBasicAuth("admin", "admin")
	req.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error sending request: %s", err)
	}

	log.Print(req)
	log.Print(resp)

}

