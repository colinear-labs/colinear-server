package walletgen

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/foxnut/go-hdwallet"
)

var mnemonic, _ = hdwallet.NewMnemonic(12, "en")
var Master, err = hdwallet.NewKey(
	hdwallet.Mnemonic(mnemonic),
)

func GenerateNewWallet(currency string) hdwallet.Wallet {
	var ctype uint32
	switch strings.ToLower(currency) {
	case "btc":
		ctype = hdwallet.BTC
	case "eth", "dai", "usdc", "usdt", "ust", "ampl":
		ctype = hdwallet.ETH
	case "bch":
		ctype = hdwallet.BCH
	case "ltc":
		ctype = hdwallet.LTC
	case "doge":
		ctype = hdwallet.DOGE
	}

	addressIndex := rand.Uint32()
	wallet, err := Master.GetWallet(hdwallet.CoinType(ctype), hdwallet.AddressIndex(addressIndex))
	if err != nil {
		fmt.Println(err)
	}
	return wallet
}
