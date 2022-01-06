package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/colinear-labs/colinear-server/flags"
	"github.com/colinear-labs/colinear-server/intents"
	"github.com/colinear-labs/colinear-server/xutil"
	"github.com/colinear-labs/colinear-server/xutil/currencies"
	"github.com/colinear-labs/colinear-server/xutil/p2pshared"

	"github.com/patrickmn/go-cache"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var Node *noise.Node = nil
var Peers = []noise.ID{}
var NodePeers = make(map[string]noise.ID)
var NodePeerCurrencies = make(map[string][]string)

func InitP2P(broadcastAddr string, port int) {

	// logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))

	// if err != nil {
	// 	panic(err)
	// }

	// defer logger.Sync()

	// port := 9871
	// broadcastIp := ipassign.GetIPv6Address()
	// broadcastIp := ipassign.GetIPv4Address()

	Node, _ = noise.NewNode(
		// noise.WithNodeLogger(logger),
		// noise.WithNodeAddress(fmt.Sprintf("[%s]:%d", broadcastIp, port)),
		noise.WithNodeAddress(broadcastAddr),
		noise.WithNodeBindPort((uint16)(port)),
	)

	// stop manually registering for now
	xutil.RegisterNodeMessages(Node)

	// Node.RegisterMessage(xutil.PeerInfo{}, xutil.UnmarshalPeerInfo)
	// Node.RegisterMessage(xutil.PaymentIntent{}, xutil.UnmarshalPaymentIntent)
	// Node.RegisterMessage(xutil.PaymentResponse{}, xutil.UnmarshalPaymentResponse)

	Node.Handle(func(ctx noise.HandlerContext) error {

		if ctx.IsRequest() {
			return nil
		}

		obj, err := ctx.DecodeMessage()
		if err != nil {
			return nil
		}

		// check between different types
		// See: https://github.com/perlin-network/noise/blob/master/example_codec_messaging_test.go#L78

		paymentResponse, ok := obj.(xutil.PaymentResponse)
		channel, _ := intents.PaymentStatusUpdateChannels.Get(paymentResponse.To)
		if ok {
			if paymentResponse.Status == xutil.Pending {
				entry, ok := intents.WatchPendingCache.Get(paymentResponse.To)
				if ok {

					intents.WatchVerifiedCache.Set(paymentResponse.To, entry, cache.DefaultExpiration)
					channel.(chan xutil.PaymentStatus) <- xutil.Pending

					fmt.Println("💸 Payment pending.")
					if flags.Mode == "single" {
						// SEND A WEBHOOK
					} else if flags.Mode == "community" {
						// SEND A WEBHOOK USING MERCHANT ID
					}
				}
			} else if paymentResponse.Status == xutil.Verified {

				intents.WatchVerifiedCache.Delete(paymentResponse.To)
				channel.(chan xutil.PaymentStatus) <- xutil.Verified

				fmt.Println("✅ Payment verified!")
				if flags.Mode == "single" {
					// SEND A WEBHOOK
				} else if flags.Mode == "community" {
					// SEND A WEBHOOK USING MERCHANT ID
				}
			}
			return nil
		}

		// PeerInfo responses

		peerInfo, ok := obj.(xutil.PeerInfo)
		if ok {

			// log occurrence
			fmt.Printf("PEERINFO %s\n", peerInfo.Marshal())

			if peerInfo.Type == xutil.Node {

				// this is embarrassingly inefficient but oh well

				added := false

			currencyCheck:
				for _, c := range currencies.Currencies {
					for _, nc := range peerInfo.Currencies {
						if c == nc {
							fmt.Printf("Full node peer supports local currency %s. Added to NodePeers.\n", c)
							NodePeers[ctx.ID().Address] = ctx.ID()
							NodePeerCurrencies[ctx.ID().Address] = peerInfo.Currencies
							added = true
							break currencyCheck
						}
					}
				}

				if !added {
					fmt.Printf("Full node peer %s does not support any local currencies.\n", ctx.ID().Address)
				}

			}
		}

		// regular byte-encoded signals

		strRes := (string)(ctx.Data())

		switch strRes {
		case "peerinfo":
			ctx.SendMessage(xutil.PeerInfo{
				Type:       xutil.Server,
				Currencies: currencies.Chains,
			})

			// add other string-formatted messages right here
		}

		return nil
	})

	k := kademlia.New()
	Node.Bind(k.Protocol())

	if err := Node.Listen(); err != nil {
		panic(err)
	}

	for _, addr := range p2pshared.BootstrapNodes {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := Node.Ping(timeoutCtx, addr+":9871")
		cancel()
		if err != nil {
			fmt.Printf("Failed to ping bootstrap node at %s.\n", addr)
		} else {
			fmt.Printf("Pinged bootstrap node at %s.\n", addr)
		}
	}

	go func() {

		for {

			Peers = k.Discover()
			fmt.Printf("%d peers detected.\n%s\n", len(Peers), fmt.Sprint(Peers))

			newNodePeers := make(map[string]noise.ID)
			newNodePeerCurrencies := make(map[string][]string)

			pingMes := []noise.ID{}

			for _, p := range Peers {
				_, ok := NodePeers[p.Address]
				if ok {
					newNodePeers[p.Address] = p
					newNodePeerCurrencies[p.Address] = NodePeerCurrencies[p.Address]
				} else {
					pingMes = append(pingMes, p)
				}
			}

			NodePeers = newNodePeers
			NodePeerCurrencies = newNodePeerCurrencies

			// Get info from new, unknown peers
			for _, p := range pingMes {
				SendGetPeerType(p.Address)
			}

			// wait 10 mins before refreshing peer list
			// Could be subject to change
			time.Sleep(10 * time.Minute)

		}

	}()

}

func SendGetPeerType(addr string) error {
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	if err := Node.Send(timeoutCtx, addr, ([]byte)("peerinfo")); err != nil {
		fmt.Printf("Error getting node type: %s\n", fmt.Sprint(err))
	} else {
		fmt.Printf("Sent GetPeerInfo to %s\n", addr)
	}

	return nil
}

func SendPaymentIntent(intent xutil.PaymentIntent) error {

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	allowableErrs := len(NodePeers)
	errCount := 0
	for _, peer := range NodePeers {

		// NOTE: probably want to check supported currencies before sending ping or something idk
		// may be expensive to do for each peer tho

		if err := Node.SendMessage(timeoutCtx, peer.Address, intent); err != nil {
			errCount += 1
			fmt.Printf("SENDMESSAGE ERROR: %s\n", err)
		} else {
			fmt.Printf("Sent intent to node %s\n", peer.Address)
		}
	}

	if allowableErrs-errCount <= 0 {
		return fmt.Errorf("not enough nodes were contacted")
	} else {
		return nil
	}

}

func SendPaymentCancellation(cancellation xutil.PaymentCancellation) error {

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	allowableErrs := len(Peers)
	errCount := 0
	for _, peer := range Peers {
		if err := Node.SendMessage(timeoutCtx, peer.Address, cancellation); err != nil {
			errCount += 1
			fmt.Printf("SENDMESSAGE ERROR: %s\n", err)
		}
	}

	if allowableErrs-errCount <= 0 {
		return fmt.Errorf("not enough nodes were contacted")
	} else {
		return nil
	}
}
