package config

import (
	"net"
	"reflect"
	"testing"

	"github.com/jrroman/caza/internal"
	"github.com/jrroman/caza/internal/testutils"
)

var (
	maskBits = 32
)

func TestValidateNetworks(t *testing.T) {
	cases := []struct {
		name          string
		networkString string
		expect        map[string]*net.IPNet
		wantError     bool
	}{
		{
			name:          "single valid network string",
			networkString: "local:127.0.0.1/32",
			expect: map[string]*net.IPNet{
				"local": testutils.CreateIPNetHelper("127.0.0.1/32"),
			},
			wantError: false,
		},
		{
			name:          "multiple valid network string",
			networkString: "local:127.0.0.1/32,router:192.168.0.0/16",
			expect: map[string]*net.IPNet{
				"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
				"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
			},
			wantError: false,
		},
		{
			name:          "network string with to many parts",
			networkString: "local::127.0.0.1/32",
			expect:        nil,
			wantError:     true,
		},
		{
			name:          "network string with invalid mask",
			networkString: "local:127.0.0.1/33",
			expect:        nil,
			wantError:     true,
		},
		{
			name:          "network string invalid cidr block",
			networkString: "local:127.0.0.10",
			expect:        nil,
			wantError:     true,
		},
		{
			name:          "missing network name",
			networkString: "127.0.0.1/32",
			expect:        nil,
			wantError:     true,
		},
	}
	for _, tc := range cases {
		got, err := validateNetworks(tc.networkString)
		if tc.wantError && err == nil {
			t.Errorf("testcase: %s expected an error", tc.name)
			continue
		}
		if !tc.wantError && err != nil {
			t.Errorf("testcase: %s did not expect error", tc.name)
			continue
		}
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("testcase: %s; got: %v, want: %v", tc.name, got, tc.expect)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	cases := []struct {
		name      string
		options   internal.Options
		expect    *Config
		wantError bool
	}{
		{
			name: "valid non cloud configuration",
			options: internal.Options{
				CloudEnabled: false,
				Networks:     "local:127.0.0.1/32,router:192.168.0.0/16",
				Region:       "",
				VpcID:        "",
			},
			expect: &Config{
				CloudEnabled: false,
				Networks: map[string]*net.IPNet{
					"local":  testutils.CreateIPNetHelper("127.0.0.1/32"),
					"router": testutils.CreateIPNetHelper("192.168.0.0/16"),
				},
				Region: "",
				VpcID:  "",
			},
			wantError: false,
		},
		{
			name: "valid cloud configuration",
			options: internal.Options{
				CloudEnabled: true,
				Networks:     "",
				Region:       "us-east-1",
				VpcID:        "abc123",
			},
			expect: &Config{
				CloudEnabled: true,
				Region:       "us-east-1",
				VpcID:        "abc123",
			},
			wantError: false,
		},
		{
			name: "no networks specified configuration",
			options: internal.Options{
				CloudEnabled: true,
				Networks:     "",
				Region:       "",
				VpcID:        "",
			},
			expect:    nil,
			wantError: true,
		},
		{
			name: "cloud configuration with no region",
			options: internal.Options{
				CloudEnabled: true,
				Region:       "",
				VpcID:        "abc123",
			},
			expect:    nil,
			wantError: true,
		},
		{
			name: "cloud configuration with no VpcID",
			options: internal.Options{
				CloudEnabled: true,
				Region:       "us-east-1",
				VpcID:        "",
			},
			expect:    nil,
			wantError: true,
		},
		{
			name: "cloud configuration with no region or VpcID",
			options: internal.Options{
				CloudEnabled: true,
				Region:       "",
				VpcID:        "",
			},
			expect:    nil,
			wantError: true,
		},
	}
	for _, tc := range cases {
		got, err := validateConfig(tc.options)
		if tc.wantError && err == nil {
			t.Errorf("testcase: %s expected an error", tc.name)
			continue
		}
		if !tc.wantError && err != nil {
			t.Errorf("testcase: %s did not expect error", tc.name)
			continue
		}
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("testcase: %s; got: %v, want: %v", tc.name, got, tc.expect)
		}
	}
}
