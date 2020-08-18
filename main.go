package main

import (
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	var g Group

	g.Add(Always, NewMonitor(context.Background()))
	g.Add(Always, NewPrinter("Chuck Bass"))
	g.Add(os.Getenv("PANICKER") == "YES", NewPanicker(5*time.Second))

	if err := g.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
