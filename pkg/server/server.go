package server

import (
	"context"
	"net"
	"sync"

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

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ebpf.LoadAndRun(eBPFEventChannel)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		ebpf.ProcessEvents(eBPFEventChannel, mergeNetworkMaps(networkSlice))
	}()
	wg.Wait()
	return nil
}
