package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"MiniProject2/Chitty_Chat"

	"google.golang.org/grpc"
)

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

	//message stream
	stream, err := client.BroadcastMessage(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	ch := clienthandle{stream: stream}
	
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl
}

//clienthandle and stream for message
type clienthandle struct {
	stream     Chitty_Chat.Chitty_Chat_BroadcastMessageClient
	clientName string
	
}

//sets name for client and status 
func (ch *clienthandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}

	ch.clientName = strings.Trim(name, "\r\n")
	ch.SendStatus()
}

//sends status to server that 
func (ch *clienthandle) SendStatus() {	

		clientMessageBox := &Chitty_Chat.BroadcastRequest{
			Name: ch.clientName,
			Message: "Has joined the Chat",
		}
	
		err := ch.stream.Send(clientMessageBox)
	
		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

}

//send message
func (ch *clienthandle) sendMessage() {

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

	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		if(mssg.Name == ""){
			fmt.Printf("%s \n", mssg.Message)
		} else {
			fmt.Printf("%s : %s \n",mssg.Name, mssg.Message)
		}
		//print message to console
	}
}
