package penta

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/ohmountain/penta/bb"
)

const (
	pentaStatusStop uint8 = iota
	pentaStatusRunning
)

type penta struct {
	mu  sync.Mutex
	prs map[string]bb.Processor[any, any]

	ctx    context.Context
	status uint8
}

func (p *penta) Insert(proccesor bb.Processor[any, any]) error {

	if p.status != pentaStatusStop {
		return errors.New("Penta is running")
	}

	p.mu.Lock()
	p.prs[proccesor.Id()] = proccesor
	p.mu.Unlock()

	return nil
}

func (p *penta) Join(proccesor bb.Processor[any, any]) error {

	if p.status != pentaStatusRunning {
		return errors.New("Penta not running")
	}

	p.mu.Lock()
	p.prs[proccesor.Id()] = proccesor
	p.mu.Unlock()

	go proccesor.Run(p.ctx)

	return nil
}

func (p *penta) Run() {

	p.mu.Lock()
	for _, prc := range p.prs {
		go prc.Run(p.ctx)
	}
	p.status = pentaStatusRunning
	p.mu.Unlock()

	if err := recover(); err != nil {
		log.Printf("遇到了一些错误: %#v", err)
	}

	for {
		select {
		case <-p.ctx.Done():
			return
		}
	}
}

func PentaWithContext(ctx context.Context) penta {
	return penta{
		mu:     sync.Mutex{},
		prs:    map[string]bb.Processor[any, any]{},
		ctx:    ctx,
		status: pentaStatusStop,
	}
}
