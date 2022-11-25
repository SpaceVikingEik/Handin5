package main

import (
	"bufio"
	"context"
	"fmt"
	Handin5 "grpcChatServer/chatserver"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	f := setLogs()
	defer f.Close()
	fmt.Println("--- Enter a Username ---")
	fmt.Printf("Your Name : ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}
	clientNameInput := strings.Trim(input, "\r\n")
	servers := []Handin5.ServicesClient{}
	for i := 0; i < 3; i++ {
		temp := "localhost:"
		serverID := temp + strconv.Itoa(5000+i)
		log.Println("Connecting : " + serverID)
		conn, err := grpc.Dial(serverID, grpc.WithInsecure())

		if err != nil {
			log.Fatalf("Failed to connect to gRPC server :: %v", err)
		}
		defer conn.Close()

		client := Handin5.NewServicesClient(conn)
		servers = append(servers, client)
	}

	reader2 := bufio.NewReader(os.Stdin)
	fmt.Println("WRITE UNIQUE ID PLS ;-;")
	input2, err := reader2.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}

	clientId, err := strconv.Atoi(input2)
	ch := clienthandle{clientName: clientNameInput, serversList: servers, clientId: clientId}

	go ch.Terminal()
	bl := make(chan bool)
	<-bl
}

type clienthandle struct {
	clientId    int
	clientName  string
	serversList []Handin5.ServicesClient
}

func (ch *clienthandle) Terminal() {
	for {
		fmt.Println("Write 'Bid' to create a bid, or 'Result' to see current highest bid")

		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Printf(" Failed to read from console :: %v", err)
		}

		input := strings.Trim(clientMessage, "\r\n")
		if input == "Bid" {
			fmt.Println("Write your bid:")

			clientBid, err := reader.ReadString('\n')
			inputBid := strings.Trim(clientBid, "\r\n")
			if err != nil {
				log.Printf(" Failed to read from console :: %v", err)
			}
			var tempResponse string
			bid, err := strconv.Atoi(inputBid)
			log.Printf("Client wants to bid: %v", bid)
			for _, element := range ch.serversList {
				bidMessage := &Handin5.BidMessage{ClientID: int64(ch.clientId), Bid: int64(bid)}
				ack, err := element.Bid(context.Background(), bidMessage)
				if err != nil {
					log.Printf(" Bid Failed %v", err)
				} else {

					tempResponse = ack.Response
					if ack.Response == "Success" {
						log.Println("Newest Highest Bid: ", int64(bid))
					}
					log.Println(tempResponse, "Server: ", element)
				}
			}
			fmt.Println(tempResponse)
		} else if input == "Result" {
			var tempResponse1 int64
			var tempResponse2 bool
			log.Println("Client requested result")
			for _, element := range ch.serversList {
				req := &Handin5.Request{ClientID: int64(ch.clientId)}
				resp, err := element.Result(context.Background(), req)
				if err != nil {
					log.Printf("Request failed %v", err)
				} else {
					tempResponse1 = resp.HighestBid
					tempResponse2 = resp.AuctionOver
					log.Println("Current Highest Bid: ", tempResponse1)
					log.Println("Auction is over: ", tempResponse2, "Server: ", element)
				}
			}
			fmt.Println(tempResponse1)
			fmt.Println("Auction is over: ", tempResponse2)
		}
	}

}

func setLogs() *os.File {
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
