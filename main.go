package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jrroman/caza/pkg/ebpf"
	"github.com/jrroman/caza/pkg/metrics"
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
	go ebpf.Run(ctx)

	select {
	case <-stopChan:
		log.Println("signal caught, shutting down")
		cancel()
	case <-ctx.Done():
		log.Println("context complete")
		return
	}
}
