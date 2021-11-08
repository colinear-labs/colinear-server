package p2p

import (
	"fmt"

	"github.com/perlin-network/noise"
)

func NoiseNode() *noise.Node {
	a, err := noise.NewNode()
	if err != nil {
		fmt.Println(err)
	}
	return a
}
