package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cross-az-analysis/pkg/ebpf"
	"cross-az-analysis/pkg/metrics"
)

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		// TODO turn metrics port into an flag
		err := metrics.Serve(":8080")
		log.Fatalf("Prometheus listener %v", err)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := ebpf.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
