package main

import (
	"../../toychain"

	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"os"
	"strings"
)

type Context struct {
	chain *toychain.Chain
	pool  Pool
}

var cliConfig = readline.Config{
	Prompt: "\033[32m>\033[0m ",
}
var flags struct {
	difficulty int
	poolSz     int
	verbose    bool
}

func init() {
	flag.IntVar(&flags.difficulty, "d", 20, "proof of work difficulty (bits)")
	flag.BoolVar(&flags.verbose, "v", false, "verbose mode")
	flag.IntVar(&flags.poolSz, "m", 0, "number of starting miners")
}

func main() {
	flag.Parse()

	chainConfig := toychain.Config{
		Difficulty: flags.difficulty,
	}
	if flags.verbose {
		chainConfig.Logger = log.New(os.Stderr, "", log.Ltime)
	}
	ctx := &Context{
		chain: toychain.NewChain(chainConfig),
	}
	if err := ctx.chain.Run(); err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < flags.poolSz; i++ {
		m := toychain.NewMiner(ctx.chain)
		ctx.pool.Add(m)
	}

	ctx.pool.RunAll()

	rl, err := readline.NewEx(&cliConfig)
	if err != nil {
		fmt.Println("readline:", err)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		switch err {
		case nil:
			argv := strings.Fields(line)
			if len(argv) == 0 {
				break
			}
			if argv[0] == "exit" {
				return
			}
			err := cliexec(ctx, argv)
			if err != nil {
				fmt.Printf("%v: %v\n", argv[0], err)
			}

		case readline.ErrInterrupt:
			if len(line) == 0 {
				return
			} else {
				break
			}

		case io.EOF:
			return

		default:
			fmt.Println("readline:", err)
		}
	}
}
