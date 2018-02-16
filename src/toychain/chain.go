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
	"bytes"
	"fmt"
	"log"
	"sync"
)

const (
	initialChainCap    = 20
	defaultBlockReward = 4
)

type Chain struct {
	sync.RWMutex
	cfg Config

	data          []*Block
	utxoset       UTXOSet
	unconfirmedTX map[string]Transaction

	blockChan chan *Block
	txChan    chan Transaction
	cfgChan   chan *Config
	stopChan  chan struct{}
}

type Config struct {
	Difficulty int
	Logger     *log.Logger
}

func (cfg *Config) Validate() error {
	if cfg.Difficulty < 1 {
		return fmt.Errorf("difficulty must be greater than 0")
	}
	return nil
}

func NewChain(cfg Config) *Chain {
	c := &Chain{
		data:          make([]*Block, 1, initialChainCap),
		utxoset:       NewUTXOSet(),
		unconfirmedTX: make(map[string]Transaction),
		blockChan:     make(chan *Block),
		txChan:        make(chan Transaction),
		cfgChan:       make(chan *Config),
		stopChan:      make(chan struct{}),
		cfg:           cfg,
	}
	c.data[0] = NewBlock()
	c.println("genesis block created")
	return c
}

func (c *Chain) Run() error {
	err := c.cfg.Validate()
	if err == nil {
		go c.loop()
	}
	return err
}

func (c *Chain) Stop() {
	c.stopChan <- struct{}{}
}

func (c *Chain) String() string {
	return fmt.Sprintf("%p", c)
}

func (c *Chain) BlockRate(i, j int) float64 {
	if i >= j || i < 0 || j < 0 {
		return 0
	}
	a := c.Block(i)
	b := c.Block(j)
	if a == nil || b == nil {
		return 0
	}
	count := b.Index - a.Index
	d := b.Timestamp.Sub(a.Timestamp).Minutes()
	return float64(count) / d
}

func (c *Chain) BlockReward() int {
	return defaultBlockReward
}

func (c *Chain) UnconfirmedTX() []Transaction {
	c.RLock()
	defer c.RUnlock()
	if len(c.unconfirmedTX) == 0 {
		return nil
	}
	s := []Transaction{}
	for _, tx := range c.unconfirmedTX {
		s = append(s, tx)
	}
	return s
}

func (c *Chain) CurrentDifficulty() int {
	return c.cfg.Difficulty
}

func (c *Chain) SetDifficulty(d int) error {
	newCfg := c.cfg
	newCfg.Difficulty = d
	err := newCfg.Validate()
	if err != nil {
		return err
	}
	c.cfgChan <- &newCfg
	c.printf("difficulty change: %v -> %v\n", c.cfg.Difficulty, d)
	return nil
}

func (c *Chain) Funds(address string) int {
	funds := 0
	for _, utxo := range c.utxoset.Get(address) {
		funds += c.fundsFromUTXO(utxo)
	}
	return funds
}

func (c *Chain) ProofIsValid(difficulty int, hash []byte) bool {
	// check first nbytes against 0
	nbytes := difficulty / 8
	for _, b := range hash[:nbytes] {
		if b != 0x00 {
			return false
		}
	}
	// check remaining bits
	nbits := uint(difficulty - 8*nbytes)
	if nbits == 0 {
		return true
	}
	if nbits > 0 {
		mask := uint(0xff >> nbits)
		b := uint(hash[nbytes])
		if (b | mask) != mask {
			return false
		}
	}
	return true
}

func (c *Chain) Block(i int) *Block {
	c.RLock()
	defer c.RUnlock()
	if i < 0 || i >= len(c.data) {
		return nil
	}
	return c.data[i]
}

func (c *Chain) BlockCount() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.data)
}

func (c *Chain) LastBlock() *Block {
	c.RLock()
	defer c.RUnlock()
	return c.data[len(c.data)-1]
}

func (c *Chain) SubmitBlock(b *Block) {
	c.blockChan <- b
}

func (c *Chain) SubmitTransaction(tx Transaction) {
	c.txChan <- tx
}

func (c *Chain) addBlock(block *Block) {
	if c.newBlockIsValid(block) == false {
		c.printf("rejecting block %x\n", block.Hash())
		return
	}
	c.printf("got block %x\n", block.Hash())

	c.Lock()
	defer c.Unlock()

	// START_UTXO OMIT
	for i, tx := range block.tx {
		// destroy inputs:
		for _, utxo := range tx.Inputs {
			addr := tx.PubKey.Address()
			c.utxoset.Remove(addr, utxo)
		}
		// distribute outputs:
		for j, output := range tx.Outputs {
			c.utxoset.Add(output, UTXO{
				BlockIndex:  block.Index,
				TxIndex:     i,
				OutputIndex: j,
			})
		}
		// ...
		// END_UTXO OMIT
		if i == 0 { // coinbase
			c.printf("executed coinbase to %v", tx.Outputs[0])
			continue
		}

		// dequeue
		digest := toString(tx.Hash())
		_, ok := c.unconfirmedTX[digest]
		if !ok {
			panic(fmt.Sprint("FIXME: unconfirmed tx not found", tx))
		}
		delete(c.unconfirmedTX, digest)
		c.printf("processed tx %x\n", tx.Hash())
	}
	c.data = append(c.data, block)
}

func (c *Chain) addUnconfirmedTx(tx Transaction) {
	c.Lock()
	defer c.Unlock()
	hash := toString(tx.Hash())
	c.unconfirmedTX[hash] = tx
	c.printf("got unconfirmed tx %x\n", tx.Hash())
}

func (c *Chain) fundsFromUTXO(utxo UTXO) int {
	c.RLock()         // OMIT
	defer c.RUnlock() //OMIT
	// ...
	defer func() { //OMIT
		if recover() != nil { //OMIT
			panic(fmt.Sprint("FIXME: invalid UTXO", utxo)) //OMIT
		} //OMIT
	}() //OMIT
	b := c.Block(utxo.BlockIndex)
	tx := b.tx[utxo.TxIndex]
	return tx.Tokens[utxo.OutputIndex]
}

func (c *Chain) newBlockIsValid(b *Block) bool {
	last := c.LastBlock()
	if !(b.Index == last.Index+1 &&
		b.Timestamp.After(last.Timestamp) &&
		bytes.Equal(b.PrevHash, last.Hash())) {
		return false
	}
	if len(b.tx) == 0 || !c.coinbaseIsValid(b.tx[0]) {
		return false
	}
	for _, tx := range b.tx[1:] {
		err := c.checkTransaction(tx)
		if err != nil {
			c.println("transaction invalid:", err)
			return false
		}
	}
	if !c.ProofIsValid(c.CurrentDifficulty(), b.Hash()) {
		return false
	}
	return true
}

func (c *Chain) coinbaseIsValid(tx Transaction) bool {
	if len(tx.Inputs) != 0 || len(tx.Outputs) != 1 {
		return false
	}
	if len(tx.Tokens) != 1 || tx.Tokens[0] != c.BlockReward() {
		return false
	}
	return true
}

func (c *Chain) checkTransaction(tx Transaction) error {
	var funds, output int
	valid := tx.PubKey.Verify(tx.Hash(), tx.Signature)
	if !valid {
		return fmt.Errorf("invalid signature")
	}
	addr := tx.PubKey.Address()
	for _, utxo := range tx.Inputs {
		if c.utxoset.Contains(addr, utxo) == false {
			return fmt.Errorf("invalid inputs")
		}
		funds += c.fundsFromUTXO(utxo)
	}
	for i := range tx.Tokens {
		output += tx.Tokens[i]
	}
	if funds < output {
		return fmt.Errorf("not enough funds: %v < %v", funds, output)
	}
	return nil
}

func (c *Chain) loop() {
	for {
		select {
		case block := <-c.blockChan:
			c.addBlock(block)

		case tx := <-c.txChan:
			c.addUnconfirmedTx(tx)

		case cfg := <-c.cfgChan:
			c.cfg = *cfg

		case <-c.stopChan:
			return
		}
	}
}

func (c *Chain) printf(format string, v ...interface{}) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.Printf(format, v...)
	}
}

func (c *Chain) println(v ...interface{}) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.Println(v...)
	}
}

func toString(data []byte) string {
	return fmt.Sprintf("%x", data)
}
