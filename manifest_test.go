package toold

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseManifest(t *testing.T) {
	m := ParseManifest("linux/amd64/node@20/yarn/go  @1.16")
	require.Equal(t, Manifest{
		OS:   Linux,
		Arch: Amd64,
		Tools: []ManifestTool{
			{
				Name:    "node",
				Version: "20",
			},
			{
				Name: "yarn",
			},
			{
				Name:    "go",
				Version: "1.16",
			},
		},
	}, m)
}
