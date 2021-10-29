package headerdata

type EthHeader struct {
	Version        uint8
	HashPrevBlock  [32]byte
	HashMerkleRoot [32]byte
	Time           uint32
	Bits           uint32
	Nonce          uint32
}
