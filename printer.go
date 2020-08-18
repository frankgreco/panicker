package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type printer struct {
	str    string
	ticker *time.Ticker
	done   chan bool
}

func NewPrinter(str string) Runnable {
	return &printer{
		str:    str,
		ticker: time.NewTicker(1 * time.Second),
		done:   make(chan bool),
	}
}

func (p printer) Run(ctx context.Context) error {
	log.Printf("starting %s printer\n", p.str)
	p.ticker.Reset(1 * time.Second)

	for {
		select {
		case <-p.done:
			return nil
		case <-p.ticker.C:
			fmt.Println(p.str)
		}
	}
}

func (p printer) Close(err error) error {
	log.Printf("closing %s printer\n", p.str)
	if p.ticker != nil {
		p.ticker.Stop()
	}
	if p.done != nil {
		p.done <- true
	}
	return nil
}
