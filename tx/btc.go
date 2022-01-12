package tx

import (
	"math/big"

	"github.com/btcsuite/btcd/wire"
)

// NOT YET DONE
func NewSignedBitcoinTx(toAddr string, amt big.Float, privateKey string) *wire.MsgTx { // figure out whether to return btcutil.Tx or wiremsg

	f, _ := amt.Float64()

	// txo := wire.NewTxOut((int64)(f*100000000), ([]byte)(""))
	// amtBtc, err := btcutil.NewAmount(f)

	tx := wire.NewMsgTx(100)                                       // figure out right version
	tx.AddTxIn(wire.NewTxIn())                                     // TXIN VS TXOUT
	tx.AddTxOut(wire.NewTxOut((int64)(f*100000000), ([]byte)(""))) // WHAT IS PKSCRIPT

	// tx, err := btcutil.NewTxFromBytes(([]byte)(""))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("tx: %v\n", tx)

	return tx
}
