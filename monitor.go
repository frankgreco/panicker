package main

import (
	"context"
	"log"
)

type process struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMonitor(ctx context.Context, cancel context.CancelFunc) Runnable {
	return &process{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *process) Run() error {
	log.Println("starting parent process monitor")
	<-p.ctx.Done()
	return nil
}

func (p *process) Close(error) error {
	log.Println("closing parent process monitor")
	p.cancel()
	return nil
}
