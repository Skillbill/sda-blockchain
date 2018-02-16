package main

import (
	crand "crypto/rand"
	sha "crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"
)

var TestLength int
var BufferSize int
var NJobs int

func init() {
	flag.IntVar(&TestLength, "s", 5, "test length (seconds)")
	flag.IntVar(&BufferSize, "size", 1024, "data size (bytes)")
	flag.IntVar(&NJobs, "j", 1, "number of jobs")
}

func main() {
	flag.Parse()
	if TestLength <= 0 || BufferSize <= 0 || NJobs <= 0 {
		flag.PrintDefaults()
		return
	}

	data := randomData(BufferSize)
	done := make(chan struct{})
	for i := 0; i < NJobs; i++ {
		go func() {
			nonce := make([]byte, binary.MaxVarintLen64)
			for v := rand.Int63(); ; v++ {
				binary.PutVarint(nonce, v)
				_ = sha.Sum256(append(data, nonce...))
				done <- struct{}{}
			}
		}()
	}

	var hashCount int
	t := time.Duration(TestLength) * time.Second
	stop := time.After(t)
	for loop := true; loop; {
		select {
		case <-done:
			hashCount++

		case <-stop:
			fmt.Println()
			loop = false
		}
	}
	fmt.Printf("%.2f kH/s\n", float64(hashCount)/1000/t.Seconds())
}

func randomData(size int) []byte {
	buf := make([]byte, size)
	r := io.LimitReader(crand.Reader, int64(size))
	r.Read(buf)
	return buf
}
