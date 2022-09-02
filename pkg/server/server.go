package server

import (
	"context"
	"log"
	"net"
	"time"

	awscloud "github.com/jrroman/caza/pkg/cloud/aws"
	"github.com/jrroman/caza/pkg/config"
	"github.com/jrroman/caza/pkg/ebpf"
)

func mergeNetworkMaps(networks []map[string]*net.IPNet) map[string]*net.IPNet {
	// if there is only one network return it
	if len(networks) == 1 {
		return networks[0]
	}
	merged := make(map[string]*net.IPNet)
	for _, nm := range networks {
		for name, network := range nm {
			merged[name] = network
		}
	}
	return merged
}

type Server struct {
	done      chan bool
	graceTime time.Duration
	startTime time.Time
}

func New() *Server {
	return &Server{
		done:      make(chan bool),
		graceTime: time.Second * 2,
		startTime: time.Now(),
	}
}

func (s *Server) Shutdown() {
	duration := time.Since(s.startTime).Seconds()
	log.Printf("caza ran for %v; server state: stopping", duration)
	s.done <- true
	// Allow time for routines to stop
	time.Sleep(s.graceTime)
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
	go ebpf.ProcessEvents(eBPFEventChannel, mergeNetworkMaps(networkSlice))
	<-ctx.Done()
	return nil
}
