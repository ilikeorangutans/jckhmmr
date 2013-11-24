package main

import (
	"log"
	"os"
	"github.com/ilikeorangutans/jckhmmr/slingclient"
	"github.com/codegangsta/cli"
	"net/url"
	"path/filepath"
	"github.com/howeyc/fsnotify"
)

var watcher fsnotify.Watcher

func main() {

	jckhmmr := &Jckhmmr{}

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
			Usage: "Password",
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
					Name: "t",
					Value: "/path/in/jcr",
					Usage: "Path in the JCR where the new file should live",
				},
			},
		},
		{
			Name: "watch",
			Usage: "Watch a directory and automatically mirror all changes",
			Action: jckhmmr.WatchDirectory,
		},
		{
			Name: "delete",
			Usage: "Deletes path",
			Action: DeletePath,
		},
	}

	app.Run(os.Args)
}

type Jckhmmr struct {

	Watcher fsnotify.Watcher

}

func DeletePath(c *cli.Context) {

	if len(c.Args()) < 1 {
		log.Panic("Not enough paths given.")
	}

	path := c.Args()[0]
	url, _ := url.Parse("http://localhost:8080")
	sc := slingclient.NewSlingClient(*url, "/", "admin", "admin")
	sc.DeletePath(path)
}

func (jckhmmr *Jckhmmr) WatchDirectory(c *cli.Context) {

	if len(c.Args()) < 1 {
		log.Panic("No directory given")
	}

	dirPath := c.Args()[0]
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Panic("Could not open directory")
	}
	defer dir.Close()

	absolutePath, _ := filepath.Abs(dir.Name())
	log.Printf("Watching %s", absolutePath)

	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()


	filepath.Walk(absolutePath, jckhmmr.WalkAndScan)

	err = Watcher.Watch(absolutePath)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()

}

func (jckhmmr *Jckhmmr) WalkAndScan(path string, info os.FileInfo, err error) error {

	if !info.IsDir() {
		return nil;
	}

	// jckhmmr.Watcher.Watch(path)
	log.Print(path)

	return nil
}

func UploadFile(c * cli.Context) {

	url, _ := url.Parse("http://localhost:8080")
	sc := slingclient.NewSlingClient(*url, "/", "admin", "admin")

	file, err := os.Open(c.String("f"))
	if err != nil {
		log.Panic("Could not open file")
	}

	path := c.String("t")

	sc.UploadFile(path, file)
}

