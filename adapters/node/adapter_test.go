package node

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/numver"
	"github.com/yankeguo/toold"
)

func TestCreateVersionExtractor(t *testing.T) {
	fn := CreateVersionExtractor(toold.AdapterOptions{
		OS:   toold.Linux,
		Arch: toold.Amd64,
	})
	v, ok := fn("node-v18.20.4-linux-x64.tar.gz")
	require.True(t, ok)
	require.Equal(t, numver.Version{18, 20, 4}, numver.Parse(v))
}
