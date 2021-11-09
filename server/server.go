// HTTP interface accessible by merchants.
//
// Contains: payment widget, REST API

package server

import (
	"fmt"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/remote/prices"
	"github.com/gofiber/fiber/v2"
)

// SUBJECT TO CHANGE
var IntentCache = make(map[string]interface{})

// Returns fiber REST API. Should be run in a goroutine with Listen()
func NewServer() *fiber.App {
	app := fiber.New()

	app.Get("/api/price", func(c *fiber.Ctx) error {
		price := prices.Price(c.Query("to"), c.Query("from"))
		if price == -1 {
			return c.SendStatus(400)
		} else {
			return c.SendString(fmt.Sprint(price))
		}
	})

	type PaymentIntentRequest struct {
		BasePrice float64 `json:basePrice`
		Currency  string  `json:currency`
		Base      string  `json:base`
	}

	app.Post("/api/createPaymentIntent", func(c *fiber.Ctx) error {

		p := new(PaymentIntentRequest)

		if err := c.BodyParser(p); err != nil {
			return c.SendStatus(400)
		}

		fmt.Println(p.Currency)
		fmt.Println(p.Base)

		price := p.BasePrice / prices.Price(p.Base, p.Currency)
		if price == -1 {
			fmt.Println("Problem with price")
			return c.SendStatus(400)
		} else {
			return c.JSON(fiber.Map{
				"amount": price,
			})
		}
	})

	// Serve static widget
	app.Static("/", "./public")

	return app
}
