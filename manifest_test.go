package toold

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseManifest(t *testing.T) {
	m := ParseManifest("linux/amd64/ww/bb/go  @1.16")
	require.Equal(t, Manifest{
		OS:   osLinux,
		Arch: archAmd64,
		Tools: map[string]string{
			"go": "1.16",
		},
	}, m)
}
