package flags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloneMode(t *testing.T) {
	cases := []struct {
		name       string
		cloneMode  string
		expected   CloneMode
		shouldFail bool
	}{
		{
			name:      "default-branch",
			cloneMode: "default-branch",
			expected:  CloneModeDefaultBranch,
		},
		{
			name:      "all-branches",
			cloneMode: "all-branches",
			expected:  CloneModeAllBranches,
		},
		{
			name:       "invalid",
			cloneMode:  "invalid",
			expected:   CloneModeInvalid, // we don't check this, just a stub
			shouldFail: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cm := CloneMode(tc.cloneMode)
			require.Equal(t, cm.String(), tc.cloneMode)

			if tc.shouldFail {
				require.Error(t, cm.Valid(), "expected to fail but didn't")
			} else {
				require.Equal(t, tc.expected, cm)
				require.NoError(t, cm.Valid(), "should not fail but did")
			}
		})
	}
}
