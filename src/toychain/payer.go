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
	"sort"
)

type Payer interface {
	PayToAddress(adress string, tokens int) error
	// ...
}

func NewPayer(chain *Chain) Payer {
	return &payer{
		chain:   chain,
		privKey: crypto.NewPrivateKey(),
	}
}

type payer struct {
	chain   *Chain
	privKey crypto.PrivateKey
}

func (p *payer) PayToAddress(address string, amount int) error {
	myaddr := p.privKey.Public().Address()
	myutxo := p.chain.utxoset.Get(myaddr)
	if len(myutxo) == 0 {
		return fmt.Errorf("no utxo for %v", myaddr)
	}
	mytokens := make([]int, len(myutxo))
	for i, utxo := range myutxo {
		mytokens[i] = p.chain.fundsFromUTXO(utxo)
	}
	inputs, change, err := p.genInputs(myutxo, mytokens, amount)
	if err != nil {
		return err
	}
	outputs, tokens := []string{address}, []int{amount}
	if change > 0 { // give me back the change, no fees in toychain
		outputs = append(outputs, myaddr)
		tokens = append(tokens, change)
	}
	tx := NewTx(inputs, outputs, tokens)
	tx.Sign(p.privKey)
	p.chain.SubmitTransaction(tx)
	return nil
}

// generate inputs from 'utxo' to satisfy 'target'
// side effect: sort tokens
func (p *payer) genInputs(utxo []UTXO, tokens []int, target int) (inputs []UTXO, change int, err error) {
	sort.Slice(utxo, func(i, j int) bool {
		return tokens[i] < tokens[j]
	})
	amount := 0
	for i := range utxo {
		amount += p.chain.fundsFromUTXO(utxo[i])
		inputs = append(inputs, utxo[i])
		if amount >= target {
			break
		}
	}
	if amount < target {
		return nil, 0, fmt.Errorf("insufficient funds (%v < %v)", amount, tokens)
	}
	change = amount - target
	return inputs, change, nil
}
