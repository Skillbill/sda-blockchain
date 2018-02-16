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
	"time"
)

type Block struct {
	// Metadata
	Timestamp time.Time
	Index     int

	// Payload
	tx []Transaction

	// Previous block's hash
	PrevHash []byte

	// Proof of work
	Difficulty int
	Nonce      int
}

func NewBlock() *Block {
	return &Block{
		Timestamp: time.Now(),
	}
}

func (b *Block) String() string {
	return fmt.Sprintf(`Block        #%v
Timestamp:   %v
Difficulty:  %v
Nonce:       %v
Hash:        %x
PrevHash:    %x
`,
		b.Index, b.Timestamp.Format(time.ANSIC),
		b.Difficulty, b.Nonce, b.Hash(), b.PrevHash)
}

func (b *Block) Hash() []byte {
	data := []byte(fmt.Sprint(b.Index, b.Timestamp, b.Difficulty, b.Nonce, b.PrevHash, b.tx))
	return crypto.Hash(data)
}
