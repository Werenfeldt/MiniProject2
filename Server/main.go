package main

import (
	"log"
	"net"
	"os"

	//"os"

	"MiniProject2/Chitty_Chat"

	"google.golang.org/grpc"
)

// // const (
// // 	port = ":8080"
// // )

// // // server is used to implement helloworld.GreeterServer.
// // type server struct {
// // 	pb.UnimplementedChitty_ChatServer
// // }

// // // SayHello implements helloworld.GreeterServer
// // func (s *server) BroadcastMessage(ctx context.Context, in *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
// // 	log.Printf("Received: %v", in.GetMessage())
// // 	return &pb.BroadcastResponse{Message: "Message receive"}, nil
// // }

// // func (is *pb.Chitty_ChatServer) BroadcastMessage(csi pb.chitty_ChatBroadcastMessageServer) error {

// // 	// ...

// // 	}

// // func (s *server) RecieveMessage(ctx context.Context, clientUniqueCode_ int){

// // 	for{
// // 		req, err := ctx.Recv();
// // 		if err != nil {
// //             log.Printf("Error reciving request from client :: %v", err)
// //             break
// // 		}else {
// // 			messageQueObject.mu.Lock()
// //             messageQueObject.MQue = append(messageQueObject.MQue, message{MessageBody: req.Body})
// //             messageQueObject.mu.Unlock()
// //             log.Printf("%v", messageQueObject.MQue[len(messageQueObject.MQue)-1])
// // 		}

// // 	}
// // 	log.Printf("Received: %v", in.GetMessage())
// // 	return &pb.BroadcastResponse{Message: "Message receive"}, nil
// // }

// // func main() {
// // 	lis, err := net.Listen("tcp", port)
// // 	if err != nil {
// // 		log.Fatalf("failed to listen: %v", err)
// // 	}
// // 	s := grpc.NewServer()
// // 	pb.RegisterChitty_ChatServer(s, &server{})
// // 	log.Printf("server listening at %v", lis.Addr())
// // 	if err := s.Serve(lis); err != nil {
// // 		log.Fatalf("failed to serve: %v", err)
// // 	}
// // }

// // //Structs
// // type message struct {
// //     //ClientName        string
// //     MessageBody       string
// //     //MessageUniqueCode int
// //     //ClientUniqueCode  int
// // }

// // type messageQue struct {
// //     MQue []message
// //     mu sync.Mutex
// // }

// // var messageQueObject = messageQue{}

// func main() {
// 	// lis, err := net.Listen("tcp", port)
// 	// if err != nil {
// 	// 	log.Fatalf("failed to listen: %v", err)
// 	// }
// 	// s := grpc.NewServer()
// 	// pb.RegisterChitty_ChatServer(s, &pb.Chitty_ChatServer{})
// 	// log.Printf("server listening at %v", lis.Addr())
// 	// if err := s.Serve(lis); err != nil {
// 	// 	log.Fatalf("failed to serve: %v", err)
// 	// }

// 	// cs := chatserver.ChatServer{}
//     // chatserver.RegisterServicesServer(grpcServer, &cs)

// 	// Port := os.Getenv("PORT")
//     // if Port == "" {
//     //     Port = "5000"// port default : 5000 if in env port is not set
//     // }

// 	Port := ":8080"

//   // initiate listener
//     listen, err := net.Listen("tcp", Port)
//     if err != nil {
//         log.Fatalf("Could not listen @ %v :: %v", Port, err)
//     }
//   	log.Println("Listening @ "+Port)

// 	  //gRPC instance server.
//   	grpcServer := grpc.NewServer()
// 	err = grpcServer.Serve(listen)
//     if err != nil {
//         log.Fatalf("Failed to start gRPC server :: %v", err)
//     }

// 	cs := Chitty_Chat.ChatServer{}
//     Chitty_Chat.RegisterChitty_ChatServer(grpcServer, &cs)
// }
func main() {
Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080" //default Port set to 5000 if PORT is not set in env
	}

	//init listener
	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()


	//register ChatService
	cs := Chitty_Chat.ChatServer{}
	Chitty_Chat.RegisterChitty_ChatServer(grpcserver,&cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}
}