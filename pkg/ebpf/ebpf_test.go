package ebpf

import (
	"net"
	"testing"

	"github.com/jrroman/caza/internal/testutils"
)

// TODO figure out unit testing for eBPF related functions

func TestCreateNetworkPairRelation(t *testing.T) {
	cases := []struct {
		name     string
		pair     *NetworkPair
		networks map[string]*net.IPNet
		expect   NetworkProximity
	}{
		{
			name: "in network pair",
			pair: newNetworkPair(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), 80, 8080),
			networks: map[string]*net.IPNet{
				"local": testutils.CreateIPNetHelper("127.0.0.1/32"),
			},
			expect: InNetwork,
		},
		{
			name: "out of network pair",
			pair: newNetworkPair(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.2"), 80, 8080),
			networks: map[string]*net.IPNet{
				"local": testutils.CreateIPNetHelper("127.0.0.1/32"),
			},
			expect: OutNetwork,
		},
		{
			name: "external network",
			pair: newNetworkPair(net.ParseIP("10.0.0.1"), net.ParseIP("127.0.0.1"), 80, 8080),
			networks: map[string]*net.IPNet{
				"local": testutils.CreateIPNetHelper("127.0.0.1/32"),
			},
			expect: ExternalNetwork,
		},
	}
	for _, tc := range cases {
		got := createNetworkPairRelation(tc.networks, tc.pair)
		if got.proximity != tc.expect {
			t.Errorf("testcase %s got: %v, want: %v", tc.name, got.proximity, tc.expect)
		}
	}
}
