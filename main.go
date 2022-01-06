package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/colinear-labs/colinear-server/config"
	"github.com/colinear-labs/colinear-server/flags"
	"github.com/colinear-labs/colinear-server/intents"
	"github.com/colinear-labs/colinear-server/p2p"
	"github.com/colinear-labs/colinear-server/server"
	"github.com/colinear-labs/colinear-server/xutil/currencies"
	"github.com/colinear-labs/colinear-server/xutil/ipassign"
)

func main() {

	mode := flag.String("mode", "single", "select single (self-hosted) mode or community mode")
	port := flag.Int("port", 80, "server port")

	flag.Parse()

	flags.Mode = *mode

	if flags.Mode == "single" {

		config.LoadSingleConfig()

		fmt.Println("Found payout addresses:")
		for k, v := range config.SingleConfig["wallets"].(map[interface{}]interface{}) {
			name := k.(string)
			currencies.Currencies = append(currencies.Currencies, name)
			curr := currencies.CurrencyData[name]
			switch curr.Type {
			case currencies.Coin:
				currencies.Chains = append(currencies.Chains, name)
			case currencies.EthToken:
				currencies.EthTokens = append(currencies.EthTokens, name)
				// add more types (e.g. BnbToken) down here in the future
			}
			fmt.Printf("\n%s:\t%s", strings.ToUpper(fmt.Sprint(k)), v)
		}

		fmt.Println("")

		fmt.Printf("Coins: %s\n", fmt.Sprint(currencies.Chains))
		fmt.Printf("Eth Tokens: %s\n", fmt.Sprint(currencies.EthTokens))

	} else if flags.Mode == "community" {
		// preflight checks for web UI API will go here

		panic("Community mode is not implemented yet.")
	}

	intents.InitIntents()

	// Try to get valid broadcast address for p2p

	p2pPort := 9871
	broadcastAddr, err := ipassign.GetNatMapping()

	if err != nil {
		fmt.Println("Failed to get a valid NAT mapping.")
		broadcastAddr, err = ipassign.GetIPv6Address()

		if err != nil {
			fmt.Println("Failed to get a valid public IPv6 address.")
			broadcastAddr, err = ipassign.GetIPv4Address()

			if err != nil {
				fmt.Println("Failed to get a valid public IPv4 address.")
				panic("Failed to get a working broadcast address.")
			} else {
				broadcastAddr = fmt.Sprintf("%s:%d", broadcastAddr, p2pPort)
			}

		} else {
			broadcastAddr = fmt.Sprintf("[%s]:%d", broadcastAddr, p2pPort)

		}

	}

	p2p.InitP2P(broadcastAddr, p2pPort)

	server := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%d", *port)))
}
