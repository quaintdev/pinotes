package main

import (
	"fmt"
	"github.com/quaintdev/pinotes/pkg/api"
	"github.com/quaintdev/pinotes/pkg/notes"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	config := notes.GetConfig()

	fs := http.FileServer(http.Dir(config.Path))
	http.HandleFunc("/add", handleAdd)
	http.Handle("/", setHeadersAndServe(fs))
	log.Println("Starting http server on", config.InterfaceIP+":"+config.HttpPort)
	go http.ListenAndServe(config.InterfaceIP+":"+config.HttpPort, nil)

	grpcServer := grpc.NewServer()
	notesService := &notes.NoteService{}
	api.RegisterNoteServiceServer(grpcServer, notesService)
	listener, err := net.Listen("tcp", config.InterfaceIP+":"+config.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting grpc server on ", config.InterfaceIP+":"+config.GrpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func setHeadersAndServe(f http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.Header().Set("Charset", "UTF-8")
		f.ServeHTTP(writer, request)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var ns notes.NoteService
	switch r.Method {
	case "GET":
		query := r.URL.Query().Get("q")
		log.Println("query received", query)
		if len(query) < 2 {
			log.Println("expecting query parameter")
			return
		}
		note := &api.Note{}
		note.Append = true
		querySlice := strings.Split(query, "!")
		if len(querySlice) == 1 {
			note.Title = notes.GetConfig().WebNotesFile
			note.Content = querySlice[0]
		} else {
			note.Title = querySlice[0]
			note.Content = querySlice[1]
		}
		if _, err := ns.SaveNote(nil, note); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		entireNote, err := ns.ViewNote(nil, note)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		w.Write([]byte(entireNote.Notes[0].GetContent()))
	default:
		fmt.Println("method not supported")
	}
}

// To build for Raspberry PI 2, use:
// GOOS=linux GOARCH=arm GOARM=7 go build pinotes.go
