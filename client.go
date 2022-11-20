package main

//Credit: https://github.com/rrrCode9/gRPC-Bidirectional-Streaming-ChatServer/blob/main/client.go
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

	fmt.Println("--- Enter a Username to Join Chat ---")
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
		conn, err := grpc.Dial(serverID, grpc.WithInsecure()) //de forsk. servere

		if err != nil {
			log.Fatalf("Failed to connect to gRPC server :: %v", err)
		}
		defer conn.Close()

		client := Handin5.NewServicesClient(conn)
		servers = append(servers, client)
	}

	//stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	reader2 := bufio.NewReader(os.Stdin)
	fmt.Printf("WRITE UNIQUE ID PLS ;-;")
	input2, err := reader2.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}

	clientId, err := strconv.Atoi(input2)
	ch := clienthandle{clientName: clientNameInput, servers: servers, clientId: int32(clientId)}

	go ch.Terminal()
	bl := make(chan bool)
	<-bl
}

type clienthandle struct {
	clientId   int32
	clientName string
	servers    []Handin5.ServicesClient
	lamport    int32
	ctx        context.Context
}

func (ch *clienthandle) Terminal() {

	for {
		fmt.Printf("Write either 1 for bid or 2 for Request")

		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		fmt.Printf(clientMessage)
		command, err := strconv.Atoi(clientMessage)
		if command == 1 {
			fmt.Printf("HEWWO")
			ch.BidMoney()
		} else if command == 2 {
			//ch.Request()
		}
	}

	//bid api call to all servers
	//accept all replies, and print óne in terminal
	//handle server crashing - use someone elses reply
}

func (ch *clienthandle) BidMoney() {
	fmt.Printf("Write Bid as a whole number:")

	reader := bufio.NewReader(os.Stdin)

	clientMessage, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}

	bid, err := strconv.Atoi(clientMessage)
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}

	for _, element := range ch.servers {
		bidMessage := &Handin5.BidMessage{ClientID: ch.clientId, Bid: int32(bid)}
		ack, err := element.Bid(ch.ctx, bidMessage)
		if err != nil {
			log.Fatalf(" Bid Failed %v", err)
		}
		log.Printf(" %v", ack)
	}

	//bid api call to all servers
	//accept all replies, and print óne in terminal
	//handle server crashing - use someone elses reply
}

/*func (ch *clienthandle) result(){
//get result
}*/

/*func (ch *clienthandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")
		ch.lamport++
		clientMessageBox := &Videobranch.FromClient{
			Name:    ch.clientName,
			Body:    clientMessage,
			Lamport: ch.lamport,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

func (ch *clienthandle) receiveMessage() {

	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in reciving message from server :: %v", err, ch.clientName)
		}
		fmt.Printf("%s : %s \n", mssg.Name, mssg.Body)
		if mssg.Lamport > ch.lamport {
			ch.lamport = mssg.Lamport
			ch.lamport++
		} else {
			ch.lamport++
		}
		fmt.Print("( Current Local Lamport Timestamp: ")
		fmt.Printf("%v %s", ch.lamport, ")")
		fmt.Println()
		fmt.Println()

	}

}*/
