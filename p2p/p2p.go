package p2p

import (
	"context"
	"fmt"
	"time"

	"xserver/flags"
	"xserver/intents"
	"xserver/xutil"
	"xserver/xutil/currencies"

	"github.com/patrickmn/go-cache"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var BootstrapNodes = [...]string{
	"10.142.24.106", // temp local address of thinkpad
}

var Node *noise.Node = nil
var XNodePeers = []noise.ID{}

func InitP2P() {

	Node, _ = noise.NewNode(noise.WithNodeBindPort(9871))

	Node.RegisterMessage(xutil.PaymentIntent{}, xutil.UnmarshalPaymentIntent)
	Node.RegisterMessage(xutil.PaymentResponse{}, xutil.UnmarshalPaymentResponse)

	// Figure out how to broadcast supported currencies
	// on peer connect
	k := kademlia.New()
	Node.Bind(k.Protocol())

	if err := Node.Listen(); err != nil {
		panic(err)
	}

	// Handle replies from other nodes
	// NOTE: Figure out if we need waigroups (sync.WaitGroup)
	Node.Handle(func(ctx noise.HandlerContext) error {
		obj, err := ctx.DecodeMessage()
		if err != nil {
			return nil
		}

		// check between different types
		// See: https://github.com/perlin-network/noise/blob/master/example_codec_messaging_test.go#L78

		paymentResponse, ok := obj.(xutil.PaymentResponse)
		if ok {
			if paymentResponse.Status == xutil.Pending {
				entry, ok := intents.WatchPendingCache.Get(paymentResponse.To)
				if ok {
					intents.WatchVerifiedCache.Set(paymentResponse.To, entry, cache.DefaultExpiration)
					fmt.Println("💸 Payment pending.")
					if flags.Mode == "single" {
						// SEND A WEBHOOK
					} else if flags.Mode == "community" {
						// SEND A WEBHOOK USING MERCHANT ID
					}
				}
			} else if paymentResponse.Status == xutil.Verified {
				intents.WatchVerifiedCache.Delete(paymentResponse.To)
				fmt.Println("✅ Payment verified!")
				if flags.Mode == "single" {
					// SEND A WEBHOOK
				} else if flags.Mode == "community" {
					// SEND A WEBHOOK USING MERCHANT ID
				}
			}
			return nil
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

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	for _, addr := range BootstrapNodes {
		if _, err := Node.Ping(timeoutCtx, addr+":9000"); err != nil {
			fmt.Printf("Failed to ping node at %s.\n", addr)
		} else {
			fmt.Printf("Pinged node at %s.\n", addr)
		}
	}

	XNodePeers = k.Discover()
	fmt.Printf("Peers: %s", fmt.Sprint(XNodePeers))

}

func SendGetNodeType(addr string) error {
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	if err := Node.Send(timeoutCtx, addr, ([]byte)("peerinfo")); err != nil {
		fmt.Printf("Error getting node type: %s\n", fmt.Sprint(err))
	} else {
		fmt.Printf("Sent GetNodeType to %s\n", addr)
	}

	return nil
}

func SendPaymentIntent(intent xutil.PaymentIntent) error {

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	allowableErrs := len(XNodePeers)
	errCount := 0
	for _, peer := range XNodePeers {
		if err := Node.SendMessage(timeoutCtx, peer.Address, intent); err != nil {
			errCount += 1
			fmt.Printf("SENDMESSAGE ERROR: %s\n", err)
		} else {
			fmt.Printf("Sent intent to node %s", peer.Address)
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

	allowableErrs := len(XNodePeers)
	errCount := 0
	for _, peer := range XNodePeers {
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
