package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Conf struct {
	Path         string
	WebNotesFile string
	Port         string
	InterfaceIP  string
}

func main() {
	configBytes, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Printf("config not present. err: %s exiting.", err)
		return
	}
	var config Conf
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Println("invalid config present. err", err)
		return
	}

	fs := http.FileServer(http.Dir(config.Path))

	http.HandleFunc("/add", handleAdd(config))
	http.Handle("/", setHeadersAndServe(fs))

	log.Println("Starting server on", config.InterfaceIP+":"+config.Port)
	http.ListenAndServe(config.InterfaceIP+":"+config.Port, nil)
}

func setHeadersAndServe(f http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.Header().Set("Charset", "UTF-8")
		f.ServeHTTP(writer, request)
	}
}

func handleAdd(config Conf) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			query := r.URL.Query().Get("q")
			log.Println("query received", query)
			if len(query) < 2 {
				log.Println("expecting query parameter")
				return
			}
			var fileToOpen string
			var note string
			querySlice := strings.Split(query, "!")
			if len(querySlice) == 1 {
				fileToOpen = config.WebNotesFile
				note = querySlice[0]
			} else {
				fileToOpen = querySlice[0]
				note = querySlice[1]
			}

			noteFilePath := path.Join(config.Path, fileToOpen+".md")
			f, err := os.OpenFile(noteFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("can't open %s, err: %s", noteFilePath, err)
				return
			}
			defer f.Close()
			_, err = f.WriteString(note + "\n\n")
			if err != nil {
				log.Println("error writing line", err)
			}
			notesContent, err := ioutil.ReadFile(noteFilePath)
			if err != nil {
				log.Printf("can't read file. err: %s", err)
			}
			w.Write(notesContent)
		default:
			fmt.Println("method not supported")
		}
	}
}

// To build for Raspberry PI 2, use:
// GOOS=linux GOARCH=arm GOARM=7 go build main.go
