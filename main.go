package main

import (
	"github.com/Super-Secret-Crypto-Kiddies/x-server/api"
)

func main() {
	// Run this in a goroutine later.
	api := api.NewApi()
	api.Listen(":3000")
}
