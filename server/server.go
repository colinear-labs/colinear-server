// HTTP interface accessible by merchants.
//
// Contains: payment widget, REST API

package server

import (
	"encoding/json"
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
	"github.com/gofiber/websocket/v2"
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
						wallets = append(wallets, currencies.EthTokensOld...)
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

		roundFactor := math.Pow(10, float64(currencies.CurrencyData[p.Currency].Decimals))
		amount := math.Round((p.BasePrice/prices.Price(p.Base, p.Currency))*roundFactor) / roundFactor
		if amount == -1 {
			return c.SendStatus(400)
		}

		intent := xutil.PaymentIntent{
			To:       address,
			Currency: p.Currency,
			Amount:   big.NewFloat(amount),
		}

		if err := p2p.SendPaymentIntent(intent); err != nil {
			fmt.Println(err)
		}
		intents.WatchPendingCache.Set(address, intent, cache.DefaultExpiration)

		statusChannel := make(chan xutil.PaymentStatus)
		intents.WatchPendingCache.Set(address, statusChannel, cache.DefaultExpiration)

		return c.JSON(fiber.Map{
			"amount":  amount,
			"address": address,
		})
	})

	type PaymentCancellationRequest struct {
		Address string
	}

	app.Post("/api/cancelPaymentIntent", func(c *fiber.Ctx) error {
		p := new(PaymentCancellationRequest)

		if err := c.BodyParser(p); err != nil {
			return c.SendStatus(400)
		}

		cancellation := xutil.PaymentCancellation{
			Address: p.Address,
		}

		if err := p2p.SendPaymentCancellation(cancellation); err != nil {
			fmt.Println(err)
			return c.SendStatus(400)
		} else {
			return c.SendStatus(200)
		}

	})

	app.Get("/ws/:toAddr", websocket.New(func(c *websocket.Conn) {
		toAddr := c.Params("toAddr")

		var statusValue string

		statusChannel, ok := intents.PaymentStatusUpdateChannels.Get(toAddr)
		if !ok {
			statusValue = "error"
			return
		}

	statusLoop:
		for {

			// write to ws
			var a interface{}
			json.Unmarshal(([]byte)(fmt.Sprintf(`{"status": "%s"}`, statusValue)), a)
			c.WriteJSON(a)

			status, ok := <-statusChannel.(chan xutil.PaymentStatus)
			if !ok {
				statusValue = "error"
			}

			statusValue, ok = map[xutil.PaymentStatus]string{
				xutil.Empty:       "empty",
				xutil.Pending:     "pending",
				xutil.Verified:    "verified",
				xutil.IntentError: "error",
			}[status]

			if !ok {
				statusValue = "error"
			}

			// check if finished state
			switch statusValue {
			case "verified", "error":
				break statusLoop
			}

		}
	}))

	// Serve static widget
	app.Static("/widget", "./widget")

	if flags.Mode == "community" {
		app.Static("/webui", "./webui")
	}

	return app
}

// DEPRECATED; DO NOT USE THIS LOL
//
// No way to delete endpoints as of right now, so could pose a memory issue
func AddPaymentStatusWssEndpoint(app *fiber.App, toAddr string, responseChannel chan xutil.PaymentStatus) {

	app.Get(fmt.Sprintf("/ws/%s", toAddr), websocket.New(func(c *websocket.Conn) {

	statusLoop:
		for {

			status := <-responseChannel

			var statusValue string
			switch status {
			case xutil.Empty:
				statusValue = "empty"
			case xutil.Pending:
				statusValue = "pending"
			case xutil.Verified:
				statusValue = "verified"
				break statusLoop
			case xutil.IntentError:
				statusValue = "error"
				break statusLoop
			default:
				statusValue = "error"
			}

			var a interface{}
			json.Unmarshal(([]byte)(fmt.Sprintf(`{"status": "%s"}`, statusValue)), a)
			c.WriteJSON(a)
		}
	}))
}
