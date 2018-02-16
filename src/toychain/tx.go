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
	"./crypto"
	"fmt"
)

type Transaction struct {
	Inputs  []UTXO
	Outputs []string // addresses
	Tokens  []int    // (not using a map just to look like bitcoin)

	PubKey    crypto.PublicKey
	Signature crypto.Signature
}

func NewCoinbase(address string, tokens int) Transaction {
	tx := makeTx(nil, []string{address}, []int{tokens})
	return tx
}

func NewTx(inputs []UTXO, outputs []string, tokens []int) Transaction {
	tx := makeTx(inputs, outputs, tokens)
	if !tx.isValid() {
		panic("invalid tx")
	}
	return tx
}

func (tx *Transaction) String() string {
	return fmt.Sprint(
		tx.Inputs,
		tx.Outputs,
		tx.Tokens,
		tx.PubKey,
	)
}

func (tx *Transaction) Hash() []byte {
	data := []byte(fmt.Sprint(tx))
	h := crypto.Hash(data)
	return h[:]
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) {
	tx.PubKey = privKey.Public()
	tx.Signature = privKey.Sign(tx.Hash())
}

func makeTx(inputs []UTXO, outputs []string, tokens []int) Transaction {
	tx := Transaction{
		Inputs:  make([]UTXO, len(inputs)),
		Outputs: make([]string, len(outputs)),
		Tokens:  make([]int, len(tokens)),
	}
	copy(tx.Inputs, inputs)
	copy(tx.Outputs, outputs)
	copy(tx.Tokens, tokens)
	return tx
}

func (tx *Transaction) isValid() bool {
	if len(tx.Outputs) < 1 || len(tx.Outputs) != len(tx.Tokens) {
		return false
	}
	for _, v := range tx.Tokens {
		if v <= 0 {
			return false
		}
	}
	// no inputs means coinbase tx
	if len(tx.Inputs) == 0 && len(tx.Outputs) != 1 {
		return false
	}
	return true
}
