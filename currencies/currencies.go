package currencies

type Currency struct {
	Name      string
	UrlPrefix string
	Decimals  uint
}

var Currencies = map[string]Currency{
	"btc": {
		Name:      "Bitcoin",
		UrlPrefix: "bitcoin",
		Decimals:  8,
	},
	"eth": {
		Name:      "Ethereum",
		UrlPrefix: "ethereum",
		Decimals:  18,
	},
	"bch": {
		Name:      "Bitcoin Cash",
		UrlPrefix: "bitcoincash",
		Decimals:  8,
	},
	"ltc": {
		Name:      "Litecoin",
		UrlPrefix: "litecoin",
		Decimals:  8,
	},
	"doge": {
		Name:      "Dogecoin",
		UrlPrefix: "dogecoin",
		Decimals:  8,
	},
	"dai": {
		Name:      "Dai",
		UrlPrefix: "dai",
		Decimals:  18,
	},
	"usdt": {
		Name:      "Tether",
		UrlPrefix: "tether",
		Decimals:  6,
	},
	"usdc": {
		Name:      "USDC",
		UrlPrefix: "usdc",
		Decimals:  18,
	},
}
