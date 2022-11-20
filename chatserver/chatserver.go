package Handin5

//Credit: https://github.com/rrrCode9/gRPC-Bidirectional-Streaming-ChatServer/blob/main/client.go
import (
	context "context"
	"sync"
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
}

func (is *ChatServer) Bid(ctx context.Context, bid *BidMessage) (*Ack, error) {
	bidM := bid
	isRegistered := false

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

}

func (is *ChatServer) Result(ctx context.Context, req *Request) (*ResultReply, error) {

	tempReply := &ResultReply{
		AuctionOver: false, //check
		HighestBid:  int64(is.currentHighestBid),
	}

	return tempReply, nil
}
