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
	Path string
	WebNotesFile string
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
	http.HandleFunc("/add", handleAdd(config))
	log.Println("Starting server on 8008")
	http.ListenAndServe(":8008", nil)
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
			querySlice := strings.Split(query, "#")
			if len(querySlice)==1{
				fileToOpen = config.WebNotesFile
				note = querySlice[0]
			}else{
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
