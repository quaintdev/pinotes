package main

import (
	"encoding/json"
	"flag"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
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
	DefaultTopic     string
	DefaultExtension string
	Port             string
	InterfaceIP      string
	CanViewTopics    bool
	DataStoreName    string
}

//Topic holds name of the topic and its content
type Topic struct {
	Id      uint
	Topic   string
	Content string
}

func main() {
	boolPtr := flag.Bool("migrate", false, "set migrate=true to do the migration")
	flag.Parse()

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

	var ds DataStore
	err = ds.Init(config.DataStoreName)
	if err != nil {
		log.Println(err)
		return
	}

	if *boolPtr {
		migrate(ds)
		return
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/add", handleBrowserRequest(ds)).Methods(http.MethodGet)
	router.HandleFunc("/topic/{topic}", handleViewTopic(ds)).Methods(http.MethodGet)
	router.HandleFunc("/topic/{topic}", handleDeleteTopic(ds)).Methods(http.MethodDelete)
	router.HandleFunc("/topic/{topic}", handleUpdateTopic(ds)).Methods(http.MethodPost)
	router.HandleFunc("/list", handleListTopics(ds)).Methods(http.MethodGet)

	log.Println("Started server on", config.InterfaceIP+":"+config.Port)
	log.Fatal(http.ListenAndServe(config.InterfaceIP+":"+config.Port, router))
}

func handleUpdateTopic(ds DataStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var topic Topic
		err := json.NewDecoder(r.Body).Decode(&topic)
		if err != nil {
			log.Println("error decoding request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		existingTopic := ds.ViewTopic(&Topic{Topic: topic.Topic})
		if len(existingTopic.Content) == 0 {
			ds.NewTopic(topic)
		} else {
			topic.Id = existingTopic.Id
			ds.UpdateTopic(topic)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleDeleteTopic(ds DataStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := mux.Vars(r)["topic"]
		err := ds.DeleteTopic(Topic{Topic: topic})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleListTopics(ds DataStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ds.ListTopics())
	}
}

func migrate(ds DataStore) {
	files, err := os.ReadDir(config.Path)
	if err != nil {
		log.Println("Unable to read directory contents. Disable migrate option.")
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			log.Println("Migrating ...", f.Name())
			file, err := os.ReadFile(path.Join(config.Path, f.Name()))
			if err != nil {
				log.Println("unable to read file", f.Name())
				return
			}
			var n Topic
			n.Topic = strings.Replace(f.Name(), config.DefaultExtension, "", -1)
			n.Content = string(file)
			if len(ds.ViewTopic(&Topic{Topic: n.Topic}).Content) == 0 {
				ds.NewTopic(n)
			} else {
				ds.UpdateTopic(n)
			}
		}
	}
	log.Println("Migration completed successfully.")
}

func handleViewTopic(ds DataStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		output := r.URL.Query().Get("output")
		topic := mux.Vars(r)["topic"]
		content := ds.ViewTopic(&Topic{Topic: topic}).Content
		if output == "text" {
			w.Write([]byte(content))
		} else {
			w.Write(markdown.ToHTML([]byte(content), nil, nil))
		}
	}
}

func handleBrowserRequest(ds DataStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var n Topic
		query := r.URL.Query().Get("q")
		log.Println("query received", query)
		if len(query) < 2 {
			log.Println("expecting query parameter")
			return
		}
		n.Topic = config.DefaultTopic
		n.Content = query

		existingTopic := ds.ViewTopic(&Topic{Topic: n.Topic})
		if len(existingTopic.Content) == 0 {
			ds.NewTopic(n)
		} else {
			var sb strings.Builder
			sb.WriteString(existingTopic.Content)
			sb.WriteString(n.Content)
			n.Content = sb.String()
			n.Id = existingTopic.Id
			ds.UpdateTopic(n)
		}
		if config.CanViewTopics {
			w.Write([]byte(n.Content))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// To build for Raspberry PI 2, use:
// GOOS=linux GOARCH=arm GOARM=7 go build github.com/quaintdev/pinotes
// For Raspberry PI Zero, use:
// GOOS=linux GOARCH=arm GOARM=5 go build github.com/quaintdev/pinotes
