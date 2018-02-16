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
	"fmt"
	"sync"
)

type UTXO struct {
	OutputIndex int
	TxIndex     int
	BlockIndex  int
}

type UTXOSet struct {
	sync.RWMutex
	data map[string][]UTXO
}

func NewUTXOSet() UTXOSet {
	return UTXOSet{
		data: make(map[string][]UTXO),
	}
}

func (set UTXOSet) Contains(address string, utxo UTXO) bool {
	s := set.Get(address)
	for i := range s {
		if s[i] == utxo {
			return true
		}
	}
	return false
}

func (set UTXOSet) Get(address string) []UTXO {
	set.RLock()
	defer set.RUnlock()

	s := set.data[address]
	if len(s) == 0 {
		return s
	}
	t := make([]UTXO, len(s))
	copy(t, s)
	return t
}

func (set UTXOSet) Add(address string, utxo UTXO) {
	set.Lock()
	defer set.Unlock()

	s := set.data[address]
	set.data[address] = append(s, utxo)
}

func (set UTXOSet) Remove(address string, utxo UTXO) {
	set.Lock()
	defer set.Unlock()

	s := set.data[address]
	if len(s) == 0 {
		fmt.Printf("warning: could not remove", utxo, "utxo not found")
		return
	}
	index := -1
	for i, v := range s {
		if v == utxo {
			index = i
			break
		}
	}
	set.data[address] = append(s[:index], s[index+1:]...)
}
