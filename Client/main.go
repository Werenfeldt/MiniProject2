package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// "os"
	"MiniProject2/Chitty_Chat"
	//"time"

	"google.golang.org/grpc"
)

// // const (
// // 	address     = "localhost:8080"
// // 	defaultName = "chittyChat"
// // )

// // func main() {
// // // Set up a connection to the server.
// // conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
// // if err != nil {
// // 	log.Fatalf("did not connect: %v", err)
// // }
// // defer conn.Close()
// // c := pd.NewChitty_ChatClient(conn)

// // //input
// // var inputMessage string

// // // Taking input from user
// // fmt.Scanln(&inputMessage)

// // // Contact the server and print out its response.
// // // Id := defaultName
// // // if len(os.Args) > 1 {
// // // 	name = os.Args[1]
// // // }
// // ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// // defer cancel()
// // r, err := c.BroadcastMessage(ctx, &pd.BroadcastRequest{Message: inputMessage})
// // if err != nil {
// // 	log.Fatalf("could not greet: %v", err)
// // }
// // log.Printf("The server: %s", r.GetMessage())
// // }

// func main() {

//     const serverID = "localhost:8080"

//     log.Println("Connecting : " + serverID)
//     conn, err := grpc.Dial(serverID, grpc.WithInsecure())

//     if err != nil {
//         log.Fatalf("Failed to connect gRPC server :: %v", err)
//     }
//     defer conn.Close()

//     client := pd.NewChitty_ChatClient(conn)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
//     stream, err := client.BroadcastMessage(ctx)
//     if err != nil {
//         log.Fatalf("Failed to get response from gRPC server :: %v", err)
//     }

//     ch := clientHandle{stream: stream}
//     ch.clientConfig()
//     go ch.sendMessage()
//     go ch.receiveMessage()

//     // block main
//     bl := make(chan bool)
//     <-bl

// }

// type clientHandle struct {
//     stream     pd.Chitty_Chat_BroadcastMessageClient
//     clientName string
// }

// func (ch *clientHandle) clientConfig() {

//     reader := bufio.NewReader(os.Stdin)
//     fmt.Printf("Your Name : ")
//     msg, err := reader.ReadString('\n')
//     if err != nil {
//         log.Fatalf("Failed to read from console :: %v", err)

//     }
//     ch.clientName = strings.TrimRight(msg, "\r\n")

// }

// func (ch *clientHandle) sendMessage() {

//     for {
//         reader := bufio.NewReader(os.Stdin)
//         clientMessage, err := reader.ReadString('\n')
//         clientMessage = strings.TrimRight(clientMessage, "\r\n")
//         if err != nil {
//             log.Printf("Failed to read from console :: %v", err)
//             continue
//         }

//         clientMessageBox := &pd.BroadcastRequest{
//             Message: clientMessage,
// 			// Name: ch.clientName,
//             // Body: clientMessage,
//         }

//         err = ch.stream.Send(clientMessageBox)

//         if err != nil {
//             log.Printf("Error while sending to server :: %v", err)
//         }

//     }
// }

// func (ch *clientHandle) receiveMessage() {

//     for {
//         resp, err := ch.stream.Recv()
//         if err != nil {
//             log.Fatalf("can not receive %v", err)
//         }
//         log.Printf("%s : %s", resp.Message)
//     }
// }

func main() {

	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	serverID = strings.Trim(serverID, "\r\n")

	log.Println("Connecting : " + serverID)

	//connect to grpc server
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	//call ChatService to create a stream
	client := Chitty_Chat.NewChitty_ChatClient(conn)

	stream, err := client.BroadcastMessage(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	// implement communication with gRPC server
	ch := clienthandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl

}

//clienthandle
type clienthandle struct {
	stream     Chitty_Chat.Chitty_Chat_BroadcastMessageClient
	clientName string
}

func (ch *clienthandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")

}

//send message
func (ch *clienthandle) sendMessage() {

	// create a loop
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &Chitty_Chat.BroadcastRequest{
			Name: ch.clientName,
			Message: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

//receive message
func (ch *clienthandle) receiveMessage() {

	//create a loop
	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		//print message to console
		fmt.Printf("%s : %s",mssg.Name, mssg.Message)
		
	}
}
