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

//Conf configuration for the note service
type Conf struct {
	Path         string
	WebNotesFile string
	Port         string
	InterfaceIP  string
	ViewNotes    bool
}

//Note holds filename of the note and its content
type Note struct {
	FileName string
	Content  string
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

	http.HandleFunc("/add", handleAdd(&config))
	http.Handle("/", setHeadersAndServe(&config, fs))

	log.Println("Starting server on", config.InterfaceIP+":"+config.Port)
	http.ListenAndServe(config.InterfaceIP+":"+config.Port, nil)
}

func setHeadersAndServe(c *Conf, f http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !c.ViewNotes {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		writer.Header().Set("Content-Type", "text/plain")
		writer.Header().Set("Charset", "UTF-8")
		f.ServeHTTP(writer, request)
	}
}

func handleAdd(config *Conf) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var n Note
		switch r.Method {
		case "POST":
			var postReq map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&postReq)
			if err != nil {
				log.Println("error decoding request")
			}
			n.FileName = postReq["title"].(string)
			n.Content = postReq["content"].(string)
			if !n.Save(config) {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
		case "GET":
			query := r.URL.Query().Get("q")
			log.Println("query received", query)
			if len(query) < 2 {
				log.Println("expecting query parameter")
				return
			}
			querySlice := strings.Split(query, "!")
			if len(querySlice) == 1 {
				n.FileName = config.WebNotesFile
				n.Content = querySlice[0]
			} else {
				n.FileName = querySlice[0]
				n.Content = querySlice[1]
			}
			//special cases that can be auto formatted in markdown
			switch n.FileName {
			case "quotes":
				n.Content = "% " + n.Content
			case "todo":
				n.Content = "1. " + n.Content
			}
			n.Save(config)
			content, err := n.Read(config)
			if err != nil {
				log.Println("error reading note", n.FileName)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(content)
		default:
			fmt.Println("method not supported")
		}
	}
}

//Save saves the note to location specified in config
func (n *Note) Save(config *Conf) bool {
	noteFilePath := path.Join(config.Path, n.FileName+".md")
	f, err := os.OpenFile(noteFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("can't open %s, err: %s", noteFilePath, err)
		return false
	}
	defer f.Close()
	_, err = f.WriteString("\n" + n.Content)
	if err != nil {
		log.Println("error writing line", err)
		return false
	}
	return true
}

func (n *Note) Read(config *Conf) ([]byte, error) {
	noteFilePath := path.Join(config.Path, n.FileName+".md")
	notesContent, err := ioutil.ReadFile(noteFilePath)
	if err != nil {
		return notesContent, fmt.Errorf("can't read file. err: %s", err)
	}
	return notesContent, nil
}

// To build for Raspberry PI 2, use:
// GOOS=linux GOARCH=arm GOARM=7 go build main.go
// For Raspberry PI Zero, use:
// GOOS=linux GOARCH=arm GOARM=5 go build main.go
