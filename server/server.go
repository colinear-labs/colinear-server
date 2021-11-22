// HTTP interface accessible by merchants.
//
// Contains: payment widget, REST API

package server

import (
	"fmt"
	"math"
	"math/big"

	"xserver/config"
	"xserver/flags"
	"xserver/intents"
	"xserver/p2p"
	"xserver/remote/prices"
	"xserver/walletgen"
	"xserver/xutil"
	"xserver/xutil/currencies"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

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
						wallets = append(wallets, currencies.EthTokens...)
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
		BasePrice float64     `json:"basePrice"`
		Currency  string      `json:"currency"`
		Base      string      `json:"base"`
		Metadata  interface{} `json:"metadata"`
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

		intent := xutil.PaymentIntent{
			To:       address,
			Currency: p.Currency,
			Amount:   big.NewFloat(amount),
		}

		p2p.SendPaymentIntent(intent)
		intents.WatchPendingCache.Set(address, intent, cache.DefaultExpiration)

		return c.JSON(fiber.Map{
			"amount":  amount,
			"address": address,
		})
	})

	// Serve static widget
	app.Static("/widget", "./widget")

	if flags.Mode == "community" {
		app.Static("/webui", "./webui")
	}

	return app
}
