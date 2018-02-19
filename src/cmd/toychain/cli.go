package main

import (
	"../../toychain"
	"fmt"
	"strconv"
)

var errf = fmt.Errorf

type Command func(*Context, []string) error

var commands = map[string]Command{
	"block": _block,
	"diff":  _diff,
	"pool":  _pool,
	"stat":  _stat,
	"pay":   _pay,
	"help":  _help,
}

func cliexec(ctx *Context, argv []string) error {
	fn, ok := commands[argv[0]]
	if !ok {
		return errf("unknown command")
	}
	return fn(ctx, argv[1:])
}

func _help(ctx *Context, args []string) error {
	fmt.Printf(`
block last		# print the last block stats
block i			# print the ith block stats

pool			# show nodes overview
pool add		# add a node to the pool
pool add n		# add n nodes to the pool
pool del i		# remove the ith node from the pool
pool del all		# remove every node from the pool

pool run|stop i		# run|stop the mining for the ith node in the pool
pool run|stop all	# run|stop the mining for every node in the pool

pay src dst amount	# make a tx of the given amount from address src to address dst

diff			# show PoW current difficulty
diff n			# change PoW difficulty to n
stat			# show toychain status

exit			# shutdown the toychain and exit
help			# you already know this one

`)
	return nil
}

func _block(ctx *Context, args []string) error {
	var b *toychain.Block

	if len(args) == 0 {
		return invalidUsage()
	}

	switch args[0] {
	case "last":
		b = ctx.chain.LastBlock()

	default:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if i < 0 {
			i += ctx.chain.BlockCount()
		}
		b = ctx.chain.Block(i)
	}

	if b == nil {
		return errf("no such block")
	}
	fmt.Print(b)
	return nil
}

func _pool(ctx *Context, args []string) error {
	var err error
	var i int

	if len(args) == 0 {
		for i, m := range ctx.pool.Data() {
			var label string
			addr := m.Address()
			amount := ctx.chain.Funds(addr)
			if !m.Running() {
				label = "stopped"
			}
			fmt.Printf("%02d %v\t%v\t%v\n", i, addr, amount, label)
		}
		return nil
	}

	switch args[0] {
	case "add":
		n := 1
		if len(args) > 1 {
			n, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
		}
		for i := 0; i < n; i++ {
			m := toychain.NewMiner(ctx.chain)
			ctx.pool.Add(m)
			fmt.Println(m.Address())
		}

	case "del":
		if len(args) < 2 {
			return invalidUsage()
		}
		switch {
		case args[1] == "all":
			ctx.pool.RemoveAll()

		default:
			i, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			err = ctx.pool.Remove(i)
		}

	case "run":
		if len(args) < 2 {
			return invalidUsage()
		}
		switch {
		case args[1] == "all":
			ctx.pool.RunAll()

		default:
			i, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			err = ctx.pool.Run(i)
		}

	case "stop":
		if len(args) < 2 {
			return invalidUsage()
		}
		switch {
		case args[1] == "all":
			ctx.pool.StopAll()

		default:
			i, err = strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			err = ctx.pool.Stop(i)
		}

	default:
		err = errf("invalid argument %q\n", args[0])
	}
	return err
}

func _diff(ctx *Context, args []string) error {
	if len(args) == 0 {
		fmt.Println("Difficulty:", ctx.chain.CurrentDifficulty())
		return nil
	}
	d, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = ctx.chain.SetDifficulty(d)
	return err

}

func _stat(ctx *Context, ignored []string) error {
	lastIdx := ctx.chain.LastBlock().Index
	i := lastIdx - 3
	if i < 0 {
		i = 0
	}
	mt, ma := ctx.pool.Count()
	fmt.Printf(`
Current difficulty:  %v
Miners:	             %v/%v
Block count:         %v
Block rate:          %.2f B/m
Unconfirmed TX:      %v

`,
		ctx.chain.CurrentDifficulty(),
		ma, mt,
		lastIdx,
		ctx.chain.BlockRate(i, lastIdx),
		len(ctx.chain.UnconfirmedTX()),
	)
	return nil
}

func _pay(ctx *Context, args []string) error {
	if len(args) < 3 {
		return invalidUsage()
	}
	src := args[0]
	dst := args[1]
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	if amount <= 0 {
		return invalidUsage()
	}
	balance := ctx.chain.Funds(src)
	if balance < amount {
		return errf("not enough funds")
	}
	node := ctx.pool.GetByAddress(src)
	if node == nil {
		return errf("no privkey for", src)
	}
	return node.PayToAddress(dst, amount)
}

func invalidUsage() error {
	return errf("invalid usage")
}
