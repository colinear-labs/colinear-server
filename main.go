package main

import (
	"fmt"
	"os"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/server"
)

func main() {
	// Run this in a goroutine later.
	port := "80"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	server := server.NewServer()
	server.Listen(fmt.Sprintf(":%s", port))
}
