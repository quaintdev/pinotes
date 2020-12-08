protoc -I api/proto --go_out=plugins=grpc:pkg/api api/proto/note.proto
GOOS=linux GOARCH=arm GOARM=7 go build pinotes.go