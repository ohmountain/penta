package penta

import (
	"context"
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

func (p *penta) Insert(proccesor bb.Processor[any, any]) {
	p.mu.Lock()
	p.prs[proccesor.Id()] = proccesor
	p.mu.Unlock()
}

func (p *penta) Run() {

	p.mu.Lock()
	for _, prc := range p.prs {
		go prc.Run(p.ctx)
	}
	p.status = pentaStatusRunning
	p.mu.Unlock()

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
