package walletgen

import (
	"fmt"
	"math/rand"
	"strings"
	"xserver/xutil/currencies"

	"github.com/foxnut/go-hdwallet"
)

var mnemonic, _ = hdwallet.NewMnemonic(12, "en")
var Master, _ = hdwallet.NewKey(
	hdwallet.Mnemonic(mnemonic),
)

func GenerateNewWallet(currency string) hdwallet.Wallet {
	var ctype uint32
	curr := strings.ToLower(currency)

	switch curr {
	case "btc":
		ctype = hdwallet.BTC
	case "eth":
		ctype = hdwallet.ETH
	case "bch":
		ctype = hdwallet.BCH
	case "ltc":
		ctype = hdwallet.LTC
	case "doge":
		ctype = hdwallet.DOGE
	default: // else - unknown
		switch currencies.CurrencyData[curr].Type {
		case currencies.EthToken:
			ctype = hdwallet.ETH
		}
		// can check other currencies down here if necessary
	}

	addressIndex := rand.Uint32()
	wallet, err := Master.GetWallet(hdwallet.CoinType(ctype), hdwallet.AddressIndex(addressIndex))
	if err != nil {
		fmt.Println(err)
	}
	return wallet
}
