package main

import (
	"github.com/Super-Secret-Crypto-Kiddies/x-server/server"
)

func main() {
	// Run this in a goroutine later.
	server := server.NewServer()
	server.Listen(":3000")
}
