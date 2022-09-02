package server

import (
	"context"
	"log"
	"net"
	"time"

	awscloud "github.com/jrroman/caza/pkg/cloud/aws"
	"github.com/jrroman/caza/pkg/config"
	"github.com/jrroman/caza/pkg/ebpf"
	"github.com/jrroman/caza/pkg/util"
)

// Server specifies some private fields which are used in managing the runtime of
// the server. Done is used to initiate shutdown of the server and gracePeriod is the
// amount of time to wait prior to a full exit of main
type Server struct {
	done        chan bool
	gracePeriod time.Duration
	startTime   time.Time
}

func New(gracefulTimeout time.Duration) *Server {
	return &Server{
		done:        make(chan bool),
		gracePeriod: gracefulTimeout,
		startTime:   time.Now(),
	}
}

func (s *Server) Shutdown() {
	duration := time.Since(s.startTime).Seconds()
	log.Printf("caza ran for %v; server state: stopping", duration)
	s.done <- true
	// Allow time for routines to stop
	time.Sleep(s.gracePeriod)
}

func (s *Server) Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		<-s.done
		cancel()
	}()
	var networkSlice []map[string]*net.IPNet
	if len(cfg.Networks) != 0 {
		networkSlice = append(networkSlice, cfg.Networks)
	}
	if cfg.CloudEnabled {
		awscc, err := awscloud.New(cfg.Region)
		if err != nil {
			return err
		}
		cloudNetworks, err := awscc.GetNetworks(cfg.VpcID)
		if err != nil {
			return err
		}
		networkSlice = append(networkSlice, cloudNetworks)
	}
	eBPFEventChannel := make(chan *ebpf.NetworkPair)
	go ebpf.LoadAndRun(ctx, eBPFEventChannel)
	go ebpf.ProcessEvents(eBPFEventChannel, util.MergeNetworkMaps(networkSlice))
	<-ctx.Done()
	return nil
}
