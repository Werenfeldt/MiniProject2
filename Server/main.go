package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"MiniProject2/Chitty_Chat"

	"google.golang.org/grpc"
)

func main() {

	LOG_FILE := "../chittyChat_log"

	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

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
	fmt.Println("Listening @ : " + Port)

	//gRPC server instance
	grpcserver := grpc.NewServer()

	//register ChatService
	cs := Chitty_Chat.ChatServer{}
	Chitty_Chat.RegisterChitty_ChatServer(grpcserver, &cs)

	//grpc listen and serve
	err = grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	} else {

	}

}
