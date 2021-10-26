package Chitty_Chat

import (
	"log"
	"math/rand"

	"sync"
	"time"
	//pb "MiniProject2/Chitty_Chat"
)

type ChatServer struct {
    
}

func (is *ChatServer) AddClient(client client){

    if(len(clientObject.CQue) == 0){
        //log.Printf("does it exist 1?: %v", Equals(client.ClientUniqueCode))
        clientObject.mu.Lock()
        clientObject.CQue = append(clientObject.CQue, client)
        clientObject.mu.Unlock()
        //log.Printf("Number of client: %v", len(clientObject.CQue))
        //log.Printf("Clien objects: %v", clientObject.CQue)
    } else {
            if Equals(client.ClientUniqueCode) {
                //log.Printf("Number of client: %v", len(clientObject.CQue))
                //log.Printf("Clien objects: %v", clientObject.CQue)
                //log.Printf("does it exist 2?: %v", Equals(client.ClientUniqueCode))
            } else {
                //log.Printf("does it exist 3?: %v", Equals(client.ClientUniqueCode))
                clientObject.mu.Lock()
                clientObject.CQue = append(clientObject.CQue, client) 
                clientObject.mu.Unlock()
                //log.Printf("Number of client: %v", len(clientObject.CQue))
                //log.Printf("Clien objects: %v", clientObject.CQue)
            }    
    }
} 

func DropClient(client int) string{
   var s string
    for i, v := range clientObject.CQue {
        
        if(v.ClientUniqueCode == client){

            s = v.clientName
            remove(i)
            break
        }  else {
            s = ""
        }
    }
    return s
}

func remove(i int){
    clientObject.CQue[i] = clientObject.CQue[len(clientObject.CQue)-1]
    clientObject.CQue = clientObject.CQue[:len(clientObject.CQue)-1]
}



func (is *ChatServer) BroadcastMessage(csi Chitty_Chat_BroadcastMessageServer) error{

    //The first time the chat sees a client, it adds it to a que
    clientUniqueCode := rand.Intn(1e3)
    client := client{
        ClientUniqueCode: clientUniqueCode,
        client: csi,  
    }

    is.AddClient(client)
    
	go RecieveMessage(csi, clientUniqueCode)
    //stream >>> client
    errch := make(chan error)
    go sendToStream(csi, clientUniqueCode, errch)

    return <-errch
} 

func AddNameToClient(name string, clientCode int){
    for i, v := range clientObject.CQue {
        if(v.ClientUniqueCode == clientCode){
            //v.clientName = name
            clientObject.CQue[i] = client{ClientUniqueCode: v.ClientUniqueCode, clientName: name, client: v.client}
            //log.Printf("%v Navn på Client", v.clientName)
        }
    }
}

func GetNameFromClient(clientCode int) string{
    var s string
    for _, v := range clientObject.CQue {
        if(v.ClientUniqueCode == clientCode){
            s = v.clientName   
        } else {
            s = ""
        }
    }
    return s
} 

func RecieveMessage(csi Chitty_Chat_BroadcastMessageServer, clientUniqueCode int){
	
	for{
		req, err := csi.Recv();
		if err != nil {
            name := DropClient(clientUniqueCode)
            log.Printf("%v has left the chat", name)
            //log.Printf("Error reciving request from client :: %v", err)
            messageQueObject.mu.Lock()
            mssg := message{ClientName: name, MessageBody: " has left the chat", ClientUniqueCode: clientUniqueCode}
            messageQueObject.MQue = append(messageQueObject.MQue, mssg)
            messageQueObject.mu.Unlock()
            break
		}else {
            messageQueObject.mu.Lock()
            messageQueObject.MQue = append(messageQueObject.MQue, message{ClientName: req.Name, MessageBody: req.Message, ClientUniqueCode: clientUniqueCode})
            
            messageQueObject.mu.Unlock()
            if(req.Name != ""){
                clientObject.mu.Lock()
                AddNameToClient(req.Name, clientUniqueCode)

                clientObject.mu.Unlock()
            }
            log.Printf("%v", messageQueObject.MQue[len(messageQueObject.MQue)-1])


            
            
            //if the status is false, it means that is a "log-on"
            // if(!req.Status){
            //     log.Printf("Status: %v", req)
            //     messageQueObject.MQue[len(messageQueObject.MQue)-1].Status = true
            // } 
		}       
	}
}

//send to stream
func sendToStream(csi Chitty_Chat_BroadcastMessageServer, clientUniqueCode int, errh chan error) {
    
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
            //status := messageQueObject.MQue[0].Status
            messageQueObject.mu.Unlock()
            
            //Checks if the status is true, which means it is a message (not a log-on)
            //if(status){
                if(len(clientObject.CQue) == 1){
                    //if the senderCode and clientcode is the same, then there is noone else (i think -need testing)
                    
                    if(!Equals(senderUniqueCode)){
                        errhh := csi.Send(&BroadcastResponse{Name: senderName4client, Message: message4client})
                            
                        if errhh != nil {
                            errh <- errhh
                        }

                    }

                        err := csi.Send(&BroadcastResponse{Name: "", Message: "You are the alone in this chat"})
    
                        if err != nil {
                            errh <- err
                        }
                        // messageQueObject.mu.Lock()
                        // mssg := message{ClientName: "", MessageBody: "You are the first one to join this chat", ClientUniqueCode: clientUniqueCode}
                        // messageQueObject.MQue = append(messageQueObject.MQue, mssg)
                        // messageQueObject.mu.Unlock()

                

                    
                } else {
                    //goes though the list of clients and sendt to all clients except itself
                    for _, csiLocal := range clientObject.CQue {
                        //log.Printf("Csi client is: %v", senderUniqueCode)
                        if(senderUniqueCode != csiLocal.ClientUniqueCode){
                            //log.Printf("local non equal sci object is: %v", csiLocal.ClientUniqueCode)
                            
                            err := csiLocal.client.Send(&BroadcastResponse{Name: senderName4client, Message: message4client})
                            
                            if err != nil {
                                errh <- err
                            }
                        } else {
                           
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
            //}
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
    client Chitty_Chat_BroadcastMessageServer
    ClientUniqueCode  int
    clientName string 
}

type message struct {
    
    ClientName        string
    MessageBody       string
    //Status            bool
    //Should probably use this. 
    //MessageUniqueCode int
    ClientUniqueCode  int
}

type messageQue struct {
    MQue []message
    mu sync.Mutex
}

type clientQue struct {
    CQue []client
    mu sync.Mutex
}

var clientObject = clientQue{}

var messageQueObject = messageQue{}
