package p2p

import (
	"context"
	"fmt"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var BoostrapNodes = [...]string{}

func NoiseNode(
	supportedCurrencies []string,
) *noise.Node {
	a, err := noise.NewNode()
	if err != nil {
		fmt.Println(err)
	}

	// Figure out how to broadcast supported currencies
	// on peer connect

	return a
}

func FindPeers(node *noise.Node) []noise.ID {
	k := kademlia.New()
	node.Bind(k.Protocol())

	if err := node.Listen(); err != nil {
		panic(err)
	}

	for _, addr := range BoostrapNodes {
		if _, err := node.Ping(context.TODO(), addr); err != nil {
			panic(err)
		}
	}

	peers := k.Discover()
	return peers

}
