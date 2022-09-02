package ebpf

import (
	"net"
	"reflect"
	"testing"

	"github.com/jrroman/caza/internal/testutils"
)

func TestMergeNetworkMaps(t *testing.T) {
	cases := []struct {
		name     string
		networks []map[string]*net.IPNet
		expect   map[string]*net.IPNet
	}{
		{
			name: "single network map",
			networks: []map[string]*net.IPNet{
				{
					"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
					"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
				},
			},
			expect: map[string]*net.IPNet{
				"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
				"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
			},
		},
		{
			name: "two network map",
			networks: []map[string]*net.IPNet{
				{
					"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
					"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
				},
				{
					"docker": testutils.CreateIPNetHelper("172.0.0.1/12"),
				},
			},
			expect: map[string]*net.IPNet{
				"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
				"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
				"docker": testutils.CreateIPNetHelper("172.0.0.1/12"),
			},
		},
		{
			name: "three network map",
			networks: []map[string]*net.IPNet{
				{
					"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
					"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
				},
				{
					"docker": testutils.CreateIPNetHelper("172.0.0.1/12"),
				},
				{
					"kubernetes": testutils.CreateIPNetHelper("10.0.0.0/12"),
				},
			},
			expect: map[string]*net.IPNet{
				"local":      testutils.CreateIPNetHelper("127.0.0.1/32"),
				"router":     testutils.CreateIPNetHelper("192.168.0.0/16"),
				"docker":     testutils.CreateIPNetHelper("172.0.0.1/12"),
				"kubernetes": testutils.CreateIPNetHelper("10.0.0.0/12"),
			},
		},
	}
	for _, tc := range cases {
		got := mergeNetworkMaps(tc.networks)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("testcase %s got: %v, want: %v", tc.name, got, tc.expect)
		}
	}
}
