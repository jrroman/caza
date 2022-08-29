package ebpf

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/jrroman/caza/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS -type event bpf fentry.c -- -I./headers

// We need to read the events being sent from our fentry.c program via ringbuffer.
// These events occur on tcp_close events and contain the src and dest ip and port.
func readLoop(ctx context.Context, rd *ringbuf.Reader, ec chan bpfEvent) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// bpfEvent is generated by bpf2go.
	var event bpfEvent
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("received signal, exiting..")
				return
			}
			log.Printf("reading from reader: %s", err)
			continue
		}

		// Parse the ringbuf event entry into a bpfEvent structure.
		if err := binary.Read(bytes.NewBuffer(record.RawSample), NativeEndian, &event); err != nil {
			log.Printf("parsing ringbuf event: %s", err)
			continue
		}
		ec <- event
	}
}

type NetworkLocale int

const (
	InNetwork NetworkLocale = iota
	OutNetwork
	ExternalNetwork
)

type NetworkMatch struct {
	Name   string
	Locale NetworkLocale
}

// First thing we have to do is identify if the address belongs to a network which
// we own. If it does then we will return the network name and check if the destination
// address also belongs to that network. If it does then we know the request stayed inside
// the network in which the request originated.
func findMatchingNetwork(networks map[string]*net.IPNet, src, dst net.IP) NetworkMatch {
	match := NetworkMatch{}
	for name, network := range networks {
		match.Name = name
		// If the src address does not live in the network continue to the next network
		if !network.Contains(src) {
			continue
		}
		// The src address belongs to the network, check if the destination address
		// also belongs to that network. If the destination address belongs to the
		// same network return "InNetwork" enum, if not return "OutNetwork enum
		if network.Contains(dst) {
			match.Locale = InNetwork
		} else {
			match.Locale = OutNetwork
		}
		return match
	}
	// If we get here it means that the src address did not belong to any of our
	// networks so we can return "ExternalNetwork" enum
	match.Name = ""
	match.Locale = ExternalNetwork
	return match
}

// We need to process the events being sent down the event channel. As a first stab
// lets create an IP map it could look something like map[net.IP]struct{in out}
func processEvents(ctx context.Context, ec chan bpfEvent, networks map[string]*net.IPNet) {
	log.Printf("%-15s %-6s -> %-15s %-6s",
		"Src addr",
		"Port",
		"Dest addr",
		"Port",
	)
	for event := range ec {
		srcAddr, dstAddr := intToIP(event.Saddr), intToIP(event.Daddr)
		srcPort, dstPort := event.Sport, event.Dport
		log.Printf("%-15s %-6d -> %-15s %-6d",
			srcAddr, srcPort, dstAddr, dstPort)
		match := findMatchingNetwork(networks, srcAddr, dstAddr)
		switch nl := match.Locale; nl {
		case InNetwork:
			metrics.InNetwork.With(prometheus.Labels{"network": match.Name}).Inc()
		case OutNetwork:
			metrics.OutNetwork.With(prometheus.Labels{"network": match.Name}).Inc()
		case ExternalNetwork:
			log.Println("A network which is not ours")
		}
	}
}

func Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Allow current process to lock memory for eBPF resources
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Printf("Remove memory lock: %v", err)
		return
	}

	// Load pre-compiled programs and maps into the kernel
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Printf("Load BPF objects: %v", err)
		return
	}
	defer objs.Close()

	link, err := link.AttachTracing(link.TracingOptions{
		Program: objs.bpfPrograms.TcpClose,
	})
	if err != nil {
		log.Printf("Attach link tracing: %v", err)
		return
	}
	defer link.Close()

	/**
	Read new tcp events from the ring buffer event data structure
	struct event {
		u16 sport;
		u16 dport;
		u32 saddr;
		u32 daddr;
	}
	*/
	rd, err := ringbuf.NewReader(objs.bpfMaps.Events)
	if err != nil {
		log.Printf("Ringbuf new reader: %v", err)
		return
	}
	defer rd.Close()

	networks, err := createNetworkMap()
	if err != nil {
		log.Printf("Create network map: %v", err)
		return
	}

	eventChan := make(chan bpfEvent)
	go readLoop(ctx, rd, eventChan)
	processEvents(ctx, eventChan, networks)
}
