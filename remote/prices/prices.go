package prices

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/imroc/req"
)

const binanceAPI = "https://api.binance.com/api/v3/ticker/price"
const coinbaseAPI = "https://api.coinbase.com/v2/prices"
const okexAPI = "https://www.okex.com/api/v1/ticker.do"
const coingeckoAPI = "https://api.coingecko.com/api/v3/simple/price"
const huobiAPI = "https://api.huobi.pro/market/detail/merged?symbol="
const krakenAPI = "https://api.kraken.com/0/public/Ticker?pair="

const chainlinkAPI = "https://api.thegraph.com/subgraphs/name/chainlink/chainlink"
const bandAPI = "https://api.thegraph.com/subgraphs/name/bandprotocol/bandprotocol"

var exchanges = map[string]string{
	"binance":   binanceAPI,
	"coinbase":  coinbaseAPI,
	"okex":      okexAPI,
	"coingecko": coingeckoAPI,
	"huobi":     huobiAPI,
	"kraken":    krakenAPI,
}

var oracles = map[string]string{
	"chainlink": chainlinkAPI,
	"band":      bandAPI,
}

// Leave out token-only exchanges until further notice
// const uniswapAPI = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap_v2"
// const dydxAPI = "https://api.dydx.exchange/v0/tickers"
// const chainlinkAPI = "https://api.thegraph.com/subgraphs/name/chainlink/chainlink"
// const bandAPI = "https://api.thegraph.com/subgraphs/name/bandprotocol/bandprotocol"

func Price(from string, to string) float64 {
	f := strings.ToUpper(from)
	t := strings.ToUpper(to)
	price := priceCoinbase(f, t)
	return price
}

func priceBinance(from string, to string) float64 {
	res, err := req.Get(binanceAPI, req.Param{
		"symbol": strings.ToUpper(to) + strings.ToUpper(from),
	})

	if err != nil {
		fmt.Println(err)
		return -1
	}

	fmt.Sprint(res)
	return 0
}

func priceCoinbase(from string, to string) float64 {
	res, err := req.Get(coinbaseAPI + "/" + strings.ToUpper(to) + "-" + strings.ToUpper(from) + "/spot")
	if err != nil {
		return -1
	}

	body := make(map[string]interface{})

	res.ToJSON(&body)

	if body["errors"] != nil {
		return -1
	}

	data := body["data"].(map[string]interface{})["amount"]
	amount, err := strconv.ParseFloat(data.(string), 64)
	if err != nil {
		return -1
	}

	return amount
}
