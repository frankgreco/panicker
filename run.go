package main

import (
	"context"
	"log"

	"github.com/oklog/run"
)

const (
	Always = true
)

type Group struct {
	group run.Group
}

type Runnable interface {
	Run(context.Context) error
	Close(error) error
}

type Process struct {
	cancel context.CancelFunc
}

func (g *Group) Add(ctx context.Context, condition bool, name string, r Runnable) {
	if !condition {
		return
	}

	g.group.Add(func() error {
		return r.Run(ctx)
	}, func(err error) {
		if err := r.Close(err); err != nil {
			log.Printf("error closing %s: %s", name, err)
		}
	})
}

func (g *Group) Run() error {
	return g.group.Run()
}

func NewMonitor(cancel context.CancelFunc) *Process {
	return &Process{
		cancel: cancel,
	}
}

func (p *Process) Run(ctx context.Context) error {
	log.Println("starting parent process monitor")
	<-ctx.Done()
	return nil
}

func (p *Process) Close(error) error {
	log.Println("closing parent process monitor")
	p.cancel()
	return nil
}
