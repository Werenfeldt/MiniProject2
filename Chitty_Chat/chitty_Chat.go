package Chitty_Chat

import (
	"log"
	"math/rand"
	"strings"

	"sync"
	"time"
)

type ChatServer struct {
}

func (is *ChatServer) AddClient(client client) {

	if len(clientObject.CQue) == 0 {
		clientObject.mu.Lock()
		clientObject.CQue = append(clientObject.CQue, client)
		clientObject.mu.Unlock()
	} else {
		if !Equals(client.ClientUniqueCode) {
			clientObject.mu.Lock()
			clientObject.CQue = append(clientObject.CQue, client)
			clientObject.mu.Unlock()
		}
	}
}

func DropClient(client int) string {
	var s string
	for i, v := range clientObject.CQue {
		if v.ClientUniqueCode == client {
			s = v.clientName
			remove(i)
			break
		} else {
			s = ""
		}
	}
	return s
}

func remove(i int) {
	clientObject.CQue[i] = clientObject.CQue[len(clientObject.CQue)-1]
	clientObject.CQue = clientObject.CQue[:len(clientObject.CQue)-1]
}

func (is *ChatServer) PublishMessage(csi Chitty_Chat_PublishMessageServer) error {

	//The first time the chat sees a client, it adds it to a que
	clientUniqueCode := rand.Intn(1e3)
	client := client{
		ClientUniqueCode: clientUniqueCode,
		client:           csi,
	}

	is.AddClient(client)

	go RecieveMessage(csi, clientUniqueCode)
	//stream >>> client
	errch := make(chan error)
	go Broadcast(csi, clientUniqueCode, errch)

	return <-errch
}

func AddNameToClient(name string, clientCode int) {
	for i, v := range clientObject.CQue {
		if v.ClientUniqueCode == clientCode {
			clientObject.CQue[i] = client{ClientUniqueCode: v.ClientUniqueCode, clientName: name, client: v.client}
		}
	}
}

func RecieveMessage(csi Chitty_Chat_PublishMessageServer, clientUniqueCode int) {

	for {
		req, err := csi.Recv()
		if err != nil {
			//logs the leaving
			name := DropClient(clientUniqueCode)
			log.Printf("%v has left the chat", name)
			//Sends message about the leaving
			messageQueObject.mu.Lock()
			mssg := message{ClientName: name, MessageBody: " has left the chat", ClientUniqueCode: clientUniqueCode}
			messageQueObject.MQue = append(messageQueObject.MQue, mssg)
			messageQueObject.mu.Unlock()
			break
		} else {
			messageQueObject.mu.Lock()
			messageQueObject.MQue = append(messageQueObject.MQue, message{ClientName: req.Name, MessageBody: req.Message, ClientUniqueCode: clientUniqueCode, Timestamp: req.Timestamp})
			messageQueObject.mu.Unlock()

			clientObject.mu.Lock()
			AddNameToClient(req.Name, clientUniqueCode)

			clientObject.mu.Unlock()
			lastMessage := messageQueObject.MQue[len(messageQueObject.MQue)-1]
			log.Printf("%v: %v at {%v}", strings.ToUpper(lastMessage.ClientName), lastMessage.MessageBody, lastMessage.Timestamp)
		}
	}
}

func Broadcast(csi Chitty_Chat_PublishMessageServer, clientUniqueCode int, errh chan error) {

	for {

		for {
			time.Sleep(500 * time.Millisecond)
			messageQueObject.mu.Lock()
			if len(messageQueObject.MQue) == 0 {
				messageQueObject.mu.Unlock()
				break
			}
			senderUniqueCode := messageQueObject.MQue[0].ClientUniqueCode
			senderName4client := messageQueObject.MQue[0].ClientName
			message4client := messageQueObject.MQue[0].MessageBody
			timestamp4client := messageQueObject.MQue[0].Timestamp
			messageQueObject.mu.Unlock()

			//Checks if the status is true, which means it is a message (not a log-on)

			if len(clientObject.CQue) == 1 {

				if !Equals(senderUniqueCode) {
					errhh := csi.Send(&PublishResponse{Name: senderName4client, Message: message4client, Timestamp: timestamp4client})

					if errhh != nil {
						errh <- errhh
					}

				}

				err := csi.Send(&PublishResponse{Name: "", Message: "You are alone in this chat"})

				if err != nil {
					errh <- err
				}

			} else {
				//goes though the list of clients and sendt to all clients except itself
				for _, csiLocal := range clientObject.CQue {

					if senderUniqueCode != csiLocal.ClientUniqueCode {
						//timestamp4client++
						err := csiLocal.client.Send(&PublishResponse{Name: senderName4client, Message: message4client, Timestamp: timestamp4client})

						if err != nil {
							errh <- err
						}
					}
				}
			}

			messageQueObject.mu.Lock()
			if len(messageQueObject.MQue) >= 2 {
				messageQueObject.MQue = messageQueObject.MQue[1:] // if send success > delete message
			} else {
				messageQueObject.MQue = []message{}
			}
			messageQueObject.mu.Unlock()
		}

		time.Sleep(1 * time.Second)
	}

}

func Equals(i int) bool {
	for _, v := range clientObject.CQue {
		if v.ClientUniqueCode == i {
			return true
		}
	}
	return false
}

//Structs
type client struct {
	client           Chitty_Chat_PublishMessageServer
	ClientUniqueCode int
	clientName       string
}

type message struct {
	ClientName       string
	MessageBody      string
	ClientUniqueCode int
	Timestamp        uint32
}

type messageQue struct {
	MQue []message
	mu   sync.Mutex
}

type clientQue struct {
	CQue []client
	mu   sync.Mutex
}

var clientObject = clientQue{}

var messageQueObject = messageQue{}
