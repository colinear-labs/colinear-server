// REST API accessible by merchants.

package api

import (
	"fmt"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/remote/prices"
	"github.com/gofiber/fiber/v2"
)

// Returns fiber REST API. Should be run in a goroutine with Listen()
func NewApi() *fiber.App {
	app := fiber.New()

	app.Get("/price", func(c *fiber.Ctx) error {
		price := prices.Price(c.Query("to"), c.Query("from"))
		if price == -1 {
			return c.SendStatus(400)
		}
		return c.SendString(fmt.Sprint(price))
	})

	// Likely in future: serve static widget from local "public" directory
	//
	// Git submodule? Manually copy over? Copy via GH Actions?
	// app.Static("/widget", "./public")

	return app
}
