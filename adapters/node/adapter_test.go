package node

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/toold"
)

func TestCreateVersionExtractor(t *testing.T) {
	fn := createVersionExtractor(toold.Linux, toold.Amd64)
	v, ok := fn("node-v18.20.4-linux-x64.tar.gz")
	require.True(t, ok)
	require.Equal(t, toold.ArbitraryVersion{18, 20, 4}, toold.ParseArbitraryVersion(v))
}
