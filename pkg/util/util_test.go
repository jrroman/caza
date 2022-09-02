package util

import (
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/jrroman/caza/internal/testutils"
)

func TestEnsureEnvironmentSet(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		setEnv    bool
		wantError bool
	}{
		{
			name:      "environment not set",
			key:       "TEST_ENV",
			setEnv:    false,
			wantError: true,
		},
		{
			name:      "environment is set",
			key:       "TEST_ENV",
			setEnv:    true,
			wantError: false,
		},
	}
	for _, tc := range cases {
		if tc.setEnv {
			os.Setenv(tc.key, "test")
		}
		err := EnsureEnvironmentSet(tc.key)
		if tc.wantError && err == nil {
			t.Errorf("testcase %s expected error", tc.name)
			continue
		}
		if !tc.wantError && err != nil {
			t.Errorf("testcase %s did not expected error", tc.name)
		}
	}
}

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
