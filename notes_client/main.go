package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dhananjayksharma/dkgosql-grpc-notes/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

// create bytes chunks for streaming
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	// TODO
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := notes.NewNotesClient(conn)

	// Define the context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// define expected flag for save
	saveCmd := flag.NewFlagSet("save", flag.ExitOnError)
	saveTitle := saveCmd.String("title", "", "Give a title to your note")
	saveId := saveCmd.String("id", "", "Give a id to your note")
	saveBody := saveCmd.String("content", "", "Type what you like to remember")

	saveLargeBody := saveCmd.Bool("l", false, "flag to upload a note broken as a stream")

	//define expected flags for load
	loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
	loadKeyword := loadCmd.String("keyword", "", "A keyword you'd like to find in your notes")

	if len(os.Args) < 3 {
		fmt.Println("expected 'save' or 'load' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "save":
		// Call the server
		// TODO
		saveCmd.Parse(os.Args[2:])
		if *saveLargeBody {
			stream, err := c.SaveLargeNote(ctx)
			if err != nil {
				log.Fatalf("Fail to create stream: %v", err)
			}
			chunks := split([]byte(*saveBody), 5)
			for _, chunk := range chunks {
				note := &notes.Note{
					Title: *saveTitle,
					Id:    *saveId,
					Body:  chunk,
				}
				if err := stream.Send(note); err != nil {
					log.Fatalf("%v.Send(%v) = %v", stream, note, err)
				}
			}
			_, err = stream.CloseAndRecv()
			if err != nil {
				log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
			}
		} else {
			bodyByte := []byte(*saveBody)

			_, err := c.Save(ctx, &notes.Note{
				Title: *saveTitle,
				Id:    *saveId,
				Body:  bodyByte,
			})

			if err != nil {
				log.Fatalf("The note could not be saved: %v", err)
			}
		}

		fmt.Printf("Your note was saved: %v%vn", *saveTitle, *saveId)
	case "load":
		loadCmd.Parse(os.Args[2:])
		// Call the server
		// TODO
		loadCmd.Parse(os.Args[2:])
		note, err := c.Load(ctx, &notes.NoteSearch{
			Keyword: *loadKeyword,
		})

		if err != nil {
			log.Fatalf("The note could not be loaded: %v", err)
		}

		fmt.Printf("%vn", note)

	default:
		fmt.Println("Expected 'save' or 'load' subcommands")
		os.Exit(1)
	}
}
