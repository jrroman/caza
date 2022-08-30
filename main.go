package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/jrroman/caza/internal"
	"github.com/jrroman/caza/pkg/config"
	"github.com/jrroman/caza/pkg/ebpf"
	"github.com/jrroman/caza/pkg/metrics"
)

var (
	opts internal.Options
)

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("Parsing command line options: %v", err)
	}
	cfg, err := config.New(opts)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		err := metrics.Serve(opts.MetricsPort)
		log.Fatalf("Prometheus listener %v", err)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := ebpf.Run(ctx, cfg); err != nil {
			log.Printf("Run %v", err)
			cancel()
		}
	}()
	select {
	case <-stopChan:
		log.Println("signal caught, shutting down")
		return
	case <-ctx.Done():
		log.Println("context complete")
		return
	}
}
