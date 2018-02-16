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
	"math/rand"
	"sync"
)

const (
	maxRounds = 200
)

type Miner interface {
	Payer

	Run()
	Running() bool
	Stop()
	Address() string
}

type miner struct {
	payer
	chain    *Chain
	privKey  crypto.PrivateKey
	stopChan chan struct{}
	running  flag
}

func NewMiner(chain *Chain) Miner {
	key := crypto.NewPrivateKey()
	m := &miner{
		chain:    chain,
		privKey:  key,
		stopChan: make(chan struct{}),
	}
	m.payer = payer{
		chain:   chain,
		privKey: key,
	}
	return m
}

func (m *miner) PayToAddress(address string, amount int) error {
	return m.payer.PayToAddress(address, amount)
}

func (m *miner) Run() {
	if m.running.Get() == false {
		go m.loop()
	}
}

func (m *miner) Stop() {
	if m.running.Get() == true {
		m.stopChan <- struct{}{}
	}
}

func (m *miner) Running() bool {
	return m.running.Get()
}

func (m *miner) Address() string {
	return m.privKey.Public().Address()
}

func (m *miner) String() string {
	return fmt.Sprintf("%s", m.Address())
}

func (m *miner) prepareNextBlock() *Block {
	last := m.chain.LastBlock()

	b := NewBlock()
	b.Index = last.Index + 1
	b.Difficulty = m.chain.CurrentDifficulty()
	b.PrevHash = last.Hash()

	coinbase := NewCoinbase(m.Address(), m.chain.BlockReward())
	tx := m.chain.UnconfirmedTX()
	b.tx = append([]Transaction{coinbase}, tx...)
	return b
}

func (m *miner) mine(output chan<- *Block) {
	block := m.prepareNextBlock()
	block.Nonce = rand.Int()

	for i := 0; i < maxRounds; i++ {
		if m.chain.ProofIsValid(block.Difficulty, block.Hash()) {
			output <- block
			return
		}
		block.Nonce++
	}
	output <- nil
}

func (m *miner) loop() {
	m.running.Set()
	for {
		ch := make(chan *Block)
		go m.mine(ch)

		select {
		case block := <-ch:
			if block != nil {
				m.chain.SubmitBlock(block)
			}

		case <-m.stopChan:
			m.running.Unset()
			return
		}
	}
}

type flag struct {
	sync.RWMutex
	value bool
}

func (f *flag) Set() {
	f.Lock()
	defer f.Unlock()
	f.value = true
}

func (f *flag) Unset() {
	f.Lock()
	defer f.Unlock()
	f.value = false
}

func (f flag) Get() bool {
	f.RLock()
	defer f.RUnlock()
	return f.value
}
