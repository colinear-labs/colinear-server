package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/config"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/flags"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/p2p"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/server"
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
			fmt.Printf("\n%s:\t%s", strings.ToUpper(fmt.Sprint(k)), v)
		}

		fmt.Println("")

	} else if flags.Mode == "community" {
		// preflight checks for web UI API will go here

		panic("Community node is not implemented yet.")
	}

	// Find peers from bootstrap nodes over noise via kademlia
	p2p.FindPeers(p2p.Node)

	server := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%d", *port)))
}
