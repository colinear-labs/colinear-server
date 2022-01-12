package tx

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewSignedEthereumTx(
	toAddr string,
	amount float64,
	privateKey ecdsa.PrivateKey,
) *types.Transaction {

	var chainId int64 = 1
	to := common.HexToAddress(toAddr)
	transaction := types.MustSignNewTx(
		&privateKey,
		types.NewLondonSigner(big.NewInt(chainId)),
		types.LegacyTx{
			Nonce: 0,
			To:    &to,
			Value: big.NewInt(int64(amount * 1e18)),
		},
	)

	return transaction
}
