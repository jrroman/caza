package util

import (
	"os"
	"testing"
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
