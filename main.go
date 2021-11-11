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

		fmt.Printf("Found payout addresses:")
		for k, v := range config.SingleConfig["addresses"].(map[interface{}]interface{}) {
			fmt.Printf("\n%s:\t%s", strings.ToUpper(fmt.Sprint(k)), v)
		}
	} else {
		// preflight checks for web UI API will go here

		panic("Community node is not implemented yet.")
	}

	server := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%d", *port)))
}
