package main

import (
	"context"
	"log"
	"time"
)

type panicker struct {
	after time.Duration
}

func NewPanicker(after time.Duration) Runnable {
	return &panicker{
		after: after,
	}
}

func (p panicker) Run(context.Context) error {
	log.Println("starting panicker")

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v\n", r)
		}
	}()

	time.Sleep(p.after)
	panic("the panicker has panicked")
	return nil
}

func (panicker) Close(err error) error {
	log.Println("closing panicker")
	return nil
}
