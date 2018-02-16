//
// This is just demo code that happens to compile and run.
//
// This software is just plain non-sense, non-idomatic, bugged,
// unpolished and unfinished.
// It's just my studying notes hacked in a pkg!
//
// This stuff is intended to be used to explain things while looking
// decent and simple on the slides.
//
package toychain

import (
	"testing"
)

var (
	zeroCfg = Config{}
	address = "LnvGRRelRQFiUw0QaIbfzKO0yW4jebBcIpEcqjIejCI"
)

func TestCoinbase(t *testing.T) {
	chain := NewChain(zeroCfg)
	tx := NewCoinbase(address, chain.BlockReward())
	if chain.coinbaseIsValid(tx) == false {
		t.Error("coinbase invalid:", tx)
	}
}
