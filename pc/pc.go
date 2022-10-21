package pc

import (
	"context"
	"fmt"
	"time"

	"github.com/ohmountain/penta/bb"
)

type pc1 struct {
	id     string
	in     chan bundle
	reader bb.OutFunc[string]
}

type pc1WriteReceipit struct {
	id        string
	timestamp int64
}

func (pr pc1WriteReceipit) Id() string {
	return pr.id
}

func (pr pc1WriteReceipit) Timestamp() int64 {
	return pr.timestamp
}

type bundle struct {
	data    string
	receipt *pc1WriteReceipit
}

func PC() pc1 {
	return pc1{
		id: "1",
	}
}

func (p pc1) Id() string {
	return p.id
}

func (p *pc1) Write(data []byte) bb.WriteReceipt {
	return pc1WriteReceipit{
		id:        fmt.Sprintf("%x", time.Now().Unix()),
		timestamp: time.Now().Unix(),
	}
}

func (p *pc1) Run(ctx context.Context) {
	if p.reader == nil {
		panic("not set reader")
	}
	for {
		select {
		case <-ctx.Done():
			return
		case bd := <-p.in:
			p.reader(bd.receipt, bd.data)
		}
	}
}
