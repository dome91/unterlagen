package common

import (
	"context"
	"sync"
)

type Shutdown struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	callbacks  []func()
}

func NewShutdown() *Shutdown {
	ctx, cancelFunc := context.WithCancel(context.Background())
	shutdown := Shutdown{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	go func() {
		<-shutdown.ctx.Done()
		for _, callback := range shutdown.callbacks {
			callback()
			shutdown.wg.Done()
		}
	}()

	return &shutdown
}

func (shutdown *Shutdown) AddCallback(callback func()) {
	shutdown.callbacks = append(shutdown.callbacks, callback)
	shutdown.wg.Add(1)
}

func (shutdown *Shutdown) Execute() {
	shutdown.cancelFunc()
	shutdown.wg.Wait()
}
