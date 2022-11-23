package main

import (
	"fmt"
	Handin5 "grpcChatServer/chatserver"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {

	f := setLog()
	defer f.Close()
	var grpcserver *grpc.Server
	for i := 0; i < 3; i++ {
		Port := i + 5000

		listen, err := net.Listen("tcp", fmt.Sprintf(":%v", Port))
		if err != nil {
			log.Fatalf("Could not listen on @ %v :: %v", Port, err)
		}
		log.Println("Listening @ : ", fmt.Sprintf(":%v", Port))

		cs := Handin5.ChatServer{}

		grpcserver = grpc.NewServer()

		Handin5.RegisterServicesServer(grpcserver, &cs)

		go startServing(grpcserver, listen)
	}

	time.Sleep(20000 * time.Millisecond)
	grpcserver.Stop()
	cha := make(chan bool)
	<-cha
}

func startServing(grpcserver *grpc.Server, listen net.Listener) {

	err := grpcserver.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}

}

//Credit: https://github.com/PatrickMatthiesen/DSYS-gRPC-template/blob/master/server/server.go

func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
