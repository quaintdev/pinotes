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
var config Conf

type Conf struct {
	Path             string
	DefaultNotesFile string
	DefaultExtension string
	Port             string
	InterfaceIP      string
	CanViewNotes     bool
}

//Note holds filename of the note and its content
type Note struct {
	Topic     string
	Content   string
	Overwrite bool
}

func main() {
	configBytes, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Printf("config not present. err: %s exiting.", err)
		return
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Println("invalid config present. err", err)
		return
	}

	fs := http.FileServer(http.Dir(config.Path))

	http.HandleFunc("/add", handleAdd)
	http.Handle("/", setHeadersAndServe(fs))

	log.Println("Starting server on", config.InterfaceIP+":"+config.Port)
	http.ListenAndServe(config.InterfaceIP+":"+config.Port, nil)
}

func setHeadersAndServe(f http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if !config.CanViewNotes {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Charset", "UTF-8")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		f.ServeHTTP(w, request)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var n Note
	switch r.Method {
	case "POST":
		var postReq map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&postReq)
		if err != nil {
			log.Println("error decoding request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		n.Topic = postReq["title"].(string)
		n.Content = postReq["content"].(string)
		n.Overwrite = postReq["overwrite"].(bool)
		if !n.Save() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "GET":
		query := r.URL.Query().Get("q")
		log.Println("query received", query)
		if len(query) < 2 {
			log.Println("expecting query parameter")
			return
		}
		topicWithContent := strings.Split(query, "!")
		if len(query) == 1 {
			n.Topic = config.DefaultNotesFile
			n.Content = topicWithContent[0]
		} else {
			n.Topic = topicWithContent[0]
			n.Content = topicWithContent[1]
		}
		n.Process()
		n.Save()
		content, err := n.Read()
		if err != nil {
			log.Println("error reading note", n.Topic)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Charset", "UTF-8")
		//Adding below headers as the pinotes webextension caches notes. We don't want it to do that.
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Write(content)
	default:
		fmt.Println("method not supported")
	}
}

func (n *Note) Process() {
	switch n.Topic {
	case "quotes":
		n.Topic = n.Topic + ".txt"
		n.Content = "% " + n.Content
	case "todo":
		n.Content = "- " + n.Content
	}
}

//Save saves the note to location specified in config
func (n *Note) Save() bool {
	fileName := n.Topic
	if !strings.Contains(fileName, ".") {
		fileName = fileName + config.DefaultExtension
	}
	noteFilePath := path.Join(config.Path, fileName)
	var writeMode int
	if n.Overwrite {
		writeMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	} else {
		writeMode = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(noteFilePath, writeMode, 0644)
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

func (n *Note) Read() ([]byte, error) {
	fileName := n.Topic
	if !strings.Contains(fileName, ".") {
		fileName = fileName + config.DefaultExtension
	}
	noteFilePath := path.Join(config.Path, fileName)
	notesContent, err := ioutil.ReadFile(noteFilePath)
	if err != nil {
		return notesContent, fmt.Errorf("can't read file. err: %s", err)
	}
	return notesContent, nil
}

// To build for Raspberry PI 2, use:
// GOOS=linux GOARCH=arm GOARM=7 go build github.com/quaintdev/pinotes
// For Raspberry PI Zero, use:
// GOOS=linux GOARCH=arm GOARM=5 go build github.com/quaintdev/pinotes
