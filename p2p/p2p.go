package p2p

import (
	"context"
	"fmt"
	"time"

	"xserver/flags"
	"xserver/intents"
	"xserver/xutil"

	"github.com/patrickmn/go-cache"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var BoostrapNodes = [...]string{
	"10.142.27.177", // temp local address of thinkpad
}

var Node *noise.Node = nil
var Peers = []noise.ID{}

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
		} else {
			// handle other message types here if necessary
		}

		return nil
	})

	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	for _, address := range BoostrapNodes {
		go func(addr string) {
			if _, err := Node.Ping(timeoutCtx, addr); err != nil {
				fmt.Printf("Failed to ping node at %s.", addr)
			}
		}(address)
	}

	Peers = k.Discover()

}
