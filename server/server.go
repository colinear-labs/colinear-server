// HTTP interface accessible by merchants.
//
// Contains: payment widget, REST API

package server

import (
	"fmt"
	"math"

	"github.com/Super-Secret-Crypto-Kiddies/x-server/config"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/currencies"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/flags"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/remote/prices"
	"github.com/Super-Secret-Crypto-Kiddies/x-server/walletgen"
	"github.com/gofiber/fiber/v2"
)

// SUBJECT TO CHANGE
var IntentCache = make(map[string]interface{})

// Returns fiber REST API. Should be run in a goroutine with Listen()
func NewServer() *fiber.App {
	app := fiber.New()

	app.Get("/api/accepting", func(c *fiber.Ctx) error {
		if flags.Mode == "single" {
			wallets := []string{}
			for k, v := range config.SingleConfig["wallets"].(map[interface{}]interface{}) {
				if v != nil && v != "" {
					wallets = append(wallets, fmt.Sprint(k))
					if k == "eth" {
						for _, token := range currencies.EthTokens {
							wallets = append(wallets, token)
						}
					}
				}
			}
			return c.JSON(wallets)
		} else {
			// Parameters will be parsed here & ignored if single mode
			return c.JSON([]string{"Community node is not implemented yet."})
		}
	})

	app.Get("/api/price", func(c *fiber.Ctx) error {
		price := prices.Price(c.Query("to"), c.Query("from"))
		if price == -1 {
			return c.SendStatus(400)
		} else {
			return c.SendString(fmt.Sprint(price))
		}
	})

	type PaymentIntentRequest struct {
		BasePrice float64     `json:basePrice`
		Currency  string      `json:currency`
		Base      string      `json:base`
		Metadata  interface{} `json:metadata`
	}

	app.Post("/api/createPaymentIntent", func(c *fiber.Ctx) error {

		p := new(PaymentIntentRequest)

		if err := c.BodyParser(p); err != nil {
			return c.SendStatus(400)
		}

		wallet := walletgen.GenerateNewWallet(p.Currency)
		address, err := wallet.GetAddress()
		if err != nil {
			return c.SendStatus(400)
		}

		roundFactor := math.Pow(10, float64(currencies.Currencies[p.Currency].Decimals))
		amount := math.Round((p.BasePrice/prices.Price(p.Base, p.Currency))*roundFactor) / roundFactor
		if amount == -1 {
			return c.SendStatus(400)
		}

		return c.JSON(fiber.Map{
			"amount":  amount,
			"address": address,
		})
	})

	// Serve static widget
	app.Static("/widget", "./widget")

	return app
}
