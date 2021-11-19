package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/xutil"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var BoostrapNodes = [...]string{
	"10.142.27.177", // temp local address of thinkpad
}

var Node = ServerNode()

func ServerNode() *noise.Node {

	node, err := noise.NewNode(noise.WithNodeBindPort(9871))
	if err != nil {
		fmt.Println(err)
	}

	node.RegisterMessage(xutil.PaymentIntent{}, xutil.UnmarshalPaymentIntent)
	node.RegisterMessage(xutil.PaymentResponse{}, xutil.UnmarshalPaymentResponse)

	// Figure out how to broadcast supported currencies
	// on peer connect

	return node
}

func FindPeers(node *noise.Node) []noise.ID {

	k := kademlia.New()
	node.Bind(k.Protocol())

	if err := node.Listen(); err != nil {
		panic(err)
	}

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	for _, address := range BoostrapNodes {
		go func(addr string) {
			if _, err := node.Ping(timeoutCtx, addr); err != nil {
				fmt.Printf("Failed to ping node at %s.", addr)
			}
		}(address)
	}

	peers := k.Discover()
	return peers

}
