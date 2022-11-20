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
	//ctx         context.Context
}

func (ch *clienthandle) Terminal() {

	for {
		fmt.Println("Write 'Bid' to create a bid, or 'Result' to see current highest bid")

		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}

		input := strings.Trim(clientMessage, "\r\n")
		if input == "Bid" {
			fmt.Println("Write your bid:")

			clientBid, err := reader.ReadString('\n')
			inputBid := strings.Trim(clientBid, "\r\n")
			if err != nil {
				log.Fatalf(" Failed to read from console :: %v", err)
			}

			bid, err := strconv.Atoi(inputBid)
			for _, element := range ch.serversList {
				bidMessage := &Handin5.BidMessage{ClientID: int64(ch.clientId), Bid: int64(bid)}
				ack, err := element.Bid(context.Background(), bidMessage)
				if err != nil {
					log.Fatalf(" Bid Failed %v", err)
				}
				log.Println(ack.Response)

				/*NOTES:
				Calling this with a bid that gets accepted, results in 1 response "success" and 2 of "Fail".
				I aktiv replik. skal hver server modtage og processe de samme requests fra clients i samme rækkefølge.
				Kan vi antage et dette betyder, at de alle har en lokal HighestBid i stedet for en global?

				- HighestBid er nu lokal variabel. Giver 3 "success".*/
			}
		} else if input == "Result" {
			for _, element := range ch.serversList {
				req := &Handin5.Request{}
				resp, err := element.Result(context.Background(), req)
				if err != nil {
					log.Fatalf("Request failed %v", err)
				}
				log.Println(resp.HighestBid)
			}
		}
	}
	//bid api call to all servers
	//accept all replies, and print óne in terminal
	//handle server crashing - use someone elses reply
}
