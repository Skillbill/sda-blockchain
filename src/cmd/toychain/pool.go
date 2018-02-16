package main

import (
	"../../toychain"
	"fmt"
)

type Pool struct {
	data    []toychain.Miner
	nactive int
}

func (p *Pool) Data() []toychain.Miner {
	return p.data
}

func (p *Pool) Get(index int) (toychain.Miner, error) {
	if !p.validIndex(index) {
		return nil, fmt.Errorf("invalid index")
	}
	return p.data[index], nil
}

func (p *Pool) GetByAddress(address string) toychain.Miner {
	for _, v := range p.data {
		if v.Address() == address {
			return v
		}
	}
	return nil
}

func (p *Pool) Count() (total, active int) {
	return len(p.data), p.nactive
}

func (p *Pool) Add(m toychain.Miner) int {
	p.data = append(p.data, m)
	return len(p.data)
}

func (p *Pool) Remove(index int) error {
	if !p.validIndex(index) {
		return fmt.Errorf("invalid index")
	}
	p.data[index].Stop()
	p.data = append(p.data[:index], p.data[index+1:]...)
	return nil
}

func (p *Pool) RemoveAll() {
	p.data = []toychain.Miner{}
}

func (p *Pool) Run(index int) error {
	if !p.validIndex(index) {
		return fmt.Errorf("invalid index")
	}
	if p.data[index].Running() == false {
		p.data[index].Run()
		p.nactive++
	}
	return nil
}

func (p *Pool) Stop(index int) error {
	if !p.validIndex(index) {
		return fmt.Errorf("invalid index")
	}
	if p.data[index].Running() {
		p.data[index].Stop()
		p.nactive--
	}
	return nil
}

func (p *Pool) RunAll() {
	for _, m := range p.data {
		m.Run()
	}
	p.nactive = len(p.data)
}

func (p *Pool) StopAll() {
	for _, m := range p.data {
		m.Stop()
	}
	p.nactive = 0
}

func (p Pool) validIndex(index int) bool {
	return index >= 0 && index < len(p.data)
}
