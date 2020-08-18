package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		firstSignal := <-c
		log.Printf("received %s signal", firstSignal.String())
		cancel()
		secondSignal := <-c
		log.Printf("received %s signal", secondSignal.String())
		os.Exit(1) // second signal. Exit directly.
	}()

	var g Group
	g.Add(Always, NewMonitor(ctx, cancel))
	g.Add(Always, NewPrinter("Chuck"))
	g.Add(Always, NewPrinter("Bass"))
	g.Add(os.Getenv("PANICKER") == "YES", NewPanicker(5*time.Second))

	if err := g.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
