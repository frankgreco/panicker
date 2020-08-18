package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type process struct {
	ctx     context.Context
	cancel  context.CancelFunc
	signals chan os.Signal
}

func NewMonitor(parent context.Context) Runnable {
	ctx, cancel := context.WithCancel(parent)
	signals := make(chan os.Signal, 2)

	return &process{
		ctx:     ctx,
		cancel:  cancel,
		signals: signals,
	}
}

func (p *process) Run() error {
	log.Println("starting parent process monitor")

	signal.Notify(p.signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		firstSignal := <-p.signals
		log.Printf("received %s signal", firstSignal.String())
		p.cancel()
		secondSignal := <-p.signals
		log.Printf("received %s signal", secondSignal.String())
		os.Exit(1) // second signal. Exit directly.
	}()

	<-p.ctx.Done()
	return nil
}

func (p *process) Close(error) error {
	log.Println("closing parent process monitor")
	p.cancel()
	return nil
}
