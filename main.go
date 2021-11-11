package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/config"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/flags"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/server"
)

func main() {

	mode := flag.String("mode", "single", "select single (self-hosted) mode or community mode")
	port := flag.Int("port", 80, "server port")

	flag.Parse()

	flags.Mode = *mode
	if *mode == "single" {

		config.LoadSingleConfig()

		supportEth := false

		// currently unused
		supportBnb := false
		supportAvax := false

		fmt.Println("Found payout addresses:")
		for k, v := range config.SingleConfig["wallets"].(map[interface{}]interface{}) {
			switch k {
			case "eth":
				supportEth = true
			case "avax":
				supportAvax = true
			}
			fmt.Printf("\n%s:\t%s", strings.ToUpper(fmt.Sprint(k)), v)
		}

		fmt.Println("")

		// Smart chain token support

		if supportEth {
			fmt.Println("Found Ethereum tokens:")
			for _, v := range config.SingleConfig["preferences"].(map[interface{}]interface{})["ethereum_tokens"].([]interface{}) {
				fmt.Printf("%s  ", strings.ToUpper(fmt.Sprint(v)))
			}
		}

		if supportBnb {
			fmt.Println("Binance is currently unsupported.")
		}
		if supportAvax {
			fmt.Println("Avalanche is currently unsupported.")
		}

	} else {
		// preflight checks for web UI API will go here

		panic("Community node is not implemented yet.")
	}

	server := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%d", *port)))
}
