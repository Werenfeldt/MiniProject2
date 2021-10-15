package Chitty_Chat

import (
	"log"
	"sync"
	"time"
	//pb "MiniProject2/Chitty_Chat"
)

// const (
// 	port = ":8080"
// )

// server is used to implement helloworld.GreeterServer.
type ChatServer struct {
}




// func (is *pb.Chitty_ChatServer) BroadcastMessage(csi pb.chitty_ChatBroadcastMessageServer) error {

// 	// ...
	
// }

func (is *ChatServer) BroadcastMessage(csi Chitty_Chat_BroadcastMessageServer) error{
	// log.Printf("Received: %v", in.GetMessage())
	// return &BroadcastResponse{Message: "Message receive"}, nil

	go RecieveMessage(csi)
    //stream >>> client
    errch := make(chan error)
    go sendToStream(csi, errch)

    return <-errch
}

func RecieveMessage(csi Chitty_Chat_BroadcastMessageServer){
	
	for{
		req, err := csi.Recv();
		if err != nil {
            log.Printf("Error reciving request from client :: %v", err)
            break
		}else {
			messageQueObject.mu.Lock()
            messageQueObject.MQue = append(messageQueObject.MQue, message{ClientName: req.Name, MessageBody: req.Message})
            messageQueObject.mu.Unlock()
            log.Printf("%v", messageQueObject.MQue[len(messageQueObject.MQue)-1])
		}
		
	}

	// log.Printf("Received: %v", in.GetMessage())
	// return &pb.BroadcastResponse{Message: "Message receive"}, nil
}

//send to stream
func sendToStream(csi Chitty_Chat_BroadcastMessageServer, errh chan error) {

    for {

        for {
            time.Sleep(500 * time.Millisecond)
            messageQueObject.mu.Lock()
            if len(messageQueObject.MQue) == 0 {
                messageQueObject.mu.Unlock()
                break
            }
            //senderUniqueCode := messageQueObject.MQue[0].ClientUniqueCode
            senderName4client := messageQueObject.MQue[0].ClientName
            message4client := messageQueObject.MQue[0].MessageBody
            messageQueObject.mu.Unlock()
            
                err := csi.Send(&BroadcastResponse{Name: senderName4client, Message: message4client})

                if err != nil {
                    errh <- err
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





//Structs
type message struct {
    ClientName        string
    MessageBody       string
    //MessageUniqueCode int
    //ClientUniqueCode  int
}

type messageQue struct {
    MQue []message
    mu sync.Mutex
}

var messageQueObject = messageQue{}