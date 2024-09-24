package toold

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseManifest(t *testing.T) {
	m := ParseManifest("linux/amd64/ww/node@20/bb/go  @1.16")
	require.Equal(t, Manifest{
		OS:   osLinux,
		Arch: archAmd64,
		Tools: []ManifestTool{
			{
				Name:    "node",
				Version: "20",
			},
			{
				Name:    "go",
				Version: "1.16",
			},
		},
	}, m)
}
