package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"xserver/config"
	"xserver/flags"
	"xserver/intents"
	"xserver/p2p"
	"xserver/server"
	"xserver/xutil/currencies"
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
	p2p.InitP2P()
	server := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%d", *port)))
}
