package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/dhananjayksharma/dkgosql-grpc-notes/notes"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// Implement the notes service (notes.NotesServer interface)
type notesServer struct {
	notes.UnimplementedNotesServer
}

// Implement the notes.NotesServer interface
func (s *notesServer) Save(ctx context.Context, n *notes.Note) (*notes.NoteSaveReply, error) {
	log.Printf("Received a note to save: %v", n.Title)
	err := notes.SaveToDisk(n, "testdata")

	if err != nil {
		return &notes.NoteSaveReply{Saved: false}, err
	}

	return &notes.NoteSaveReply{Saved: true}, nil
}

// Implement the notes.NotesServer interface
func (s *notesServer) Load(ctx context.Context, search *notes.NoteSearch) (*notes.Note, error) {
	log.Printf("Received a note to load: %v", search.Keyword)
	n, err := notes.LoadFromDisk(search.Keyword, "testdata")

	if err != nil {
		return &notes.Note{}, err
	}

	return n, nil
}

// Implement the new SaveLargeNote function in the server
func (s *notesServer) SaveLargeNote(stream notes.Notes_SaveLargeNoteServer) error {
	var finalBody []byte
	var finalTitle string
	var id string
	for {
		// Get a packet
		note, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Received a note to save: %v", finalTitle)
			err := notes.SaveToDisk(&notes.Note{
				Title: finalTitle,
				Id:    id,
				Body:  finalBody,
			}, "testdata")

			if err != nil {
				stream.SendAndClose(&notes.NoteSaveReply{Saved: false})
				return err
			}

			stream.SendAndClose(&notes.NoteSaveReply{Saved: true})
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Received a chunk of the note to save: %v", note.Body)
		// Concat packet to create final note
		finalBody = append(finalBody, note.Body...)
		finalTitle = note.Title
		id = note.Id
	}
}

func main() {
	// parse arguments from the command line
	// this lets us define the port for the server
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// Check for errors
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Instantiate the server
	s := grpc.NewServer()
	// Register server method (actions the server will do)
	notes.RegisterNotesServer(s, &notesServer{})
	// Register server method (actions the server will do)
	// TODO

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
