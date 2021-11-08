// HTTP interface accessible by merchants.
//
// Contains: payment widget, REST API

package server

import (
	"fmt"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/remote/prices"
	"github.com/gofiber/fiber/v2"
)

// Returns fiber REST API. Should be run in a goroutine with Listen()
func NewServer() *fiber.App {
	app := fiber.New()

	app.Get("/api/price", func(c *fiber.Ctx) error {
		price := prices.Price(c.Query("to"), c.Query("from"))
		if price == -1 {
			return c.SendStatus(400)
		}
		return c.SendString(fmt.Sprint(price))
	})

	// Serve static widget
	app.Static("/", "./public")

	return app
}
