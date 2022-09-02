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
	"github.com/jrroman/caza/pkg/metrics"
	"github.com/jrroman/caza/pkg/server"
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
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signalCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := metrics.Serve(opts.MetricsPort)
		log.Fatalf("Prometheus listener %v", err)
	}()
	srv := server.New(cfg.GraceTime)
	go func() {
		if err := srv.Run(ctx, cfg); err != nil {
			log.Printf("Run %v", err)
			cancel()
		}
	}()
	select {
	case <-ctx.Done():
		log.Println("context done, shutting down")
		return
	case s := <-signalCh:
		log.Printf("signal %v caught, starting shutdown...", s)
		srv.Shutdown()
		cancel()
		return
	}
}
