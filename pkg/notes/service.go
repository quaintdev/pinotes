package notes

import (
	"context"
	"github.com/quaintdev/pinotes/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type NoteService struct {
}

func (ns *NoteService) SaveNote(ctx context.Context, note *api.Note) (*api.Response, error) {
	noteFilePath := path.Join(GetConfig().Path, note.Title+".md")
	var flag int
	if note.Append == true {
		flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flag = os.O_CREATE | os.O_WRONLY
	}
	f, err := os.OpenFile(noteFilePath, flag, 0644)
	if err != nil {
		log.Printf("can't open %s, err: %s", noteFilePath, err)
		return &api.Response{Success: false}, status.Errorf(codes.Internal, "error opening file")
	}
	defer f.Close()
	_, err = f.WriteString(note.Content + "\n\n")
	if err != nil {
		log.Println("error writing line", err)
		return &api.Response{Success: false}, status.Errorf(codes.Internal, "error adding note")
	}
	return &api.Response{Success: true}, nil
}

func (ns *NoteService) ViewNote(ctx context.Context, note *api.Note) (r *api.Response, err error) {
	noteFilePath := path.Join(GetConfig().Path, note.Title+".md")
	notesContent, err := ioutil.ReadFile(noteFilePath)
	if err != nil {
		log.Printf("can't read file. err: %s", err)
		return &api.Response{}, status.Errorf(codes.Internal, "error viewing notes")
	}
	r = &api.Response{
		Success: true,
		Notes: []*api.Note{
			&api.Note{
				Title:   note.Title,
				Content: string(notesContent),
			},
		},
	}
	return r, nil
}

func (ns *NoteService) DeleteNote(ctx context.Context, note *api.Note) (r *api.Response, err error) {
	noteFilePath := path.Join(GetConfig().Path, note.Title+".md")
	err = os.Remove(noteFilePath)
	if err != nil {
		return &api.Response{}, status.Errorf(codes.Internal, "error deleting notes")
	}
	return &api.Response{Success: true}, err
}

func (ns *NoteService) ListNotes(ctx context.Context, note *api.Note) (r *api.Response, err error) {
	files, err := ioutil.ReadDir(GetConfig().Path)
	if err != nil {
		log.Fatal(err)
	}
	r = &api.Response{
		Notes:   []*api.Note{},
		Success: true,
	}
	for _, f := range files {
		r.Notes = append(r.Notes, &api.Note{Title: strings.Replace(f.Name(), ".md", "", -1)})
		log.Println(f.Name())
	}
	return
}

/*
grpcurl -d '{"Title": "today", "Content": "make bananas", "user":"Rohan"}' -plaintext -import-path ./api/proto/ -proto note.proto localhost:9009 api.NoteService.SaveNote
*/
