package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"time"
)

type Header struct {
	version    int // always 1
	difficulty int
	date       Timestamp // YYMMDDhhmmss
	resource   string    // recipient
	extension  string    // unused
	random     string    // base64 encoded
	nonce      Nonce     // base64 encoded
}

func (h *Header) String() string {
	return fmt.Sprintf("%v:%v:%v:%v::%v:%v",
		h.version,
		h.difficulty, h.date,
		h.resource, h.random, h.nonce)
}

func (h *Header) Hash() [sha1.Size]byte {
	str := fmt.Sprintf("%v", h)
	return sha1.Sum([]byte(str))
}

func (h *Header) HashIsValid() bool {
	hash := h.Hash()

	// check first nbytes against 0
	nbytes := h.difficulty / 8

	for _, b := range hash[:nbytes] {
		if b != 0x00 {
			return false
		}
	}
	// check remaining bits
	nbits := h.difficulty - 8*nbytes
	if nbits > 0 {
		mask := uint(0xff) >> uint(nbits)
		b := uint(hash[nbytes])
		if (b | mask) != mask {
			return false
		}
	}
	return true
}

type Timestamp time.Time

func (ts Timestamp) String() string {
	t := time.Time(ts)
	yr, mon, day := t.Date()
	return fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		yr-2000, mon, day, t.Hour(), t.Minute(), t.Second())
}

type Nonce uint32

const (
	NonceByteSize = 3
	NonceMax      = 1<<20 - 1
)

func NewNonce() Nonce {
	buf := random(4)
	v := binary.LittleEndian.Uint32(buf) & NonceMax
	return Nonce(v)
}

func (c *Nonce) Next() error {
	*c++
	if *c > NonceMax {
		return fmt.Errorf("overflow")
	}
	return nil
}

func (c Nonce) String() string {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(c))
	return fmt.Sprintf("%v", b64(buf[:NonceByteSize]))
}

func Hashcash(difficulty int, resource string, output chan<- *Header) {
	h := &Header{
		version:    1,
		difficulty: difficulty,
		date:       Timestamp(time.Now()),
		resource:   resource,
	}
	
setup:
	h.random = b64(random(10))
	h.nonce = NewNonce()

	for {
		if h.HashIsValid() {
			output <- h
			return
		}
		err := h.nonce.Next()
		if err != nil { // overflow
			goto setup
		}
	}
}

var cfg struct {
	rsc        string
	difficulty int
	njobs      int
	verbose    bool
	help       bool
}

func init() {
	flag.StringVar(&cfg.rsc, "r", "sergio@localhost", "resource")
	flag.IntVar(&cfg.difficulty, "d", 20, "difficulty (bits)")
	flag.IntVar(&cfg.njobs, "j", 1, "number of concurrent jobs")
	flag.BoolVar(&cfg.help, "h", false, "help")
}

func main() {
	flag.Parse()
	if cfg.help || cfg.difficulty < 0 || cfg.njobs <= 0 {
		flag.PrintDefaults()
		return
	}
	ch := make(chan *Header)
	for i := 0; i < cfg.njobs; i++ {
		go Hashcash(cfg.difficulty, cfg.rsc, ch)
	}
	clock := time.Tick(100 * time.Millisecond)
	for {
		select {
		case header := <-ch:
			fmt.Println()
			fmt.Println(header)
			return

		case <-clock:
			fmt.Printf(".")
		}
	}
}

func random(size int64) []byte {
	buf := make([]byte, size)
	r := io.LimitReader(rand.Reader, size)
	r.Read(buf)
	return buf
}

func b64(buf []byte) string {
	enc := base64.StdEncoding.WithPadding(base64.NoPadding)
	return enc.EncodeToString(buf)
}
