package main

import (
	"context"
	"log"
	"fmt"

	// "os"
	pd "MiniProject2/Chitty_Chat"
	"time"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:8080"
	defaultName = "chittyChat"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pd.NewChitty_ChatClient(conn)

	//input
	var inputMessage string
  
    // Taking input from user
    fmt.Scanln(&inputMessage)

	// Contact the server and print out its response.
	// Id := defaultName
	// if len(os.Args) > 1 {
	// 	name = os.Args[1]
	// }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BroadcastMessage(ctx, &pd.BroadcastRequest{Message: inputMessage})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("The server: %s", r.GetMessage())
}