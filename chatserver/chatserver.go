package Handin5

//Credit: https://github.com/rrrCode9/gRPC-Bidirectional-Streaming-ChatServer/blob/main/client.go
import (
	context "context"
	"log"
	"sync"
	"time"
)

type chatserviceHandle struct {
	ClientBidMap map[int]clienthandle
	lo           sync.Mutex
	clientList   []clienthandle
}

type clienthandle struct {
	cName string
	id    int
}

var chatserviceHandleObject = chatserviceHandle{ClientBidMap: make(map[int]clienthandle), clientList: make([]clienthandle, 0)}

type ChatServer struct {
	UnimplementedServicesServer
	currentHighestBid int
	bidLock           sync.Mutex
	auctionOver       bool
	isNotFirstTime    bool
}

func (is *ChatServer) Bid(ctx context.Context, bid *BidMessage) (*Ack, error) {
	if !is.isNotFirstTime {
		log.Printf("Auction has begun for server!")
		go is.Timer()
		is.isNotFirstTime = true
	} else {
		time.Sleep(2000 * time.Millisecond)
	}
	bidM := bid
	isRegistered := false
	if is.auctionOver == false {
		chatserviceHandleObject.lo.Lock()
		for _, element := range chatserviceHandleObject.clientList {
			if bidM.ClientID == int64(element.id) {
				isRegistered = true
			}
		}
		if isRegistered {
			is.bidLock.Lock()
			if bidM.Bid > int64(is.currentHighestBid) {
				is.currentHighestBid = int(bidM.Bid)

				is.bidLock.Unlock()
				chatserviceHandleObject.lo.Unlock()
				return &Ack{Response: "Success"}, nil

			} else {

				is.bidLock.Unlock()
				chatserviceHandleObject.lo.Unlock()
				return &Ack{Response: "Fail"}, nil
			}

		} else {
			temp := clienthandle{cName: "", id: int(bidM.ClientID)}
			chatserviceHandleObject.clientList = append(chatserviceHandleObject.clientList, temp)
			chatserviceHandleObject.ClientBidMap[0] = temp

			is.bidLock.Lock()
			if bidM.Bid > int64(is.currentHighestBid) {
				is.currentHighestBid = int(bidM.Bid)

				is.bidLock.Unlock()
				chatserviceHandleObject.lo.Unlock()
				return &Ack{Response: "Success"}, nil

			} else {

				is.bidLock.Unlock()
				chatserviceHandleObject.lo.Unlock()
				return &Ack{Response: "Fail"}, nil
			}
		}
	} else {
		return &Ack{Response: "AuctionOver"}, nil
	}
}

func (is *ChatServer) Result(ctx context.Context, req *Request) (*ResultReply, error) {
	if !is.auctionOver {
		is.bidLock.Lock()
		tempReply := &ResultReply{
			AuctionOver: false,
			HighestBid:  int64(is.currentHighestBid),
		}
		is.bidLock.Unlock()
		return tempReply, nil
	} else {
		tempReply := &ResultReply{
			AuctionOver: true,
			HighestBid:  int64(is.currentHighestBid),
		}
		return tempReply, nil
	}
}

func (is *ChatServer) Timer() {
	timer := time.NewTimer(40000 * time.Millisecond)
	buffer1 := make(chan bool)
	go func() {
		<-timer.C
		buffer1 <- true
	}()
	var temp = <-buffer1
	is.bidLock.Lock()
	is.auctionOver = temp
	is.bidLock.Unlock()
}
