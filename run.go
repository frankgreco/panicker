package main

import (
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
	Run() error
	Close(error) error
}

func (g *Group) Add(condition bool, r Runnable) {
	if !condition {
		return
	}

	g.group.Add(func() error {
		return r.Run()
	}, func(err error) {
		if err := r.Close(err); err != nil {
			log.Printf("error closing runnable: %s", err)
		}
	})
}

func (g *Group) Run() error {
	return g.group.Run()
}
