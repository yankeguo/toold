package maven

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/toold"
)

func TestCreateVersionExtractor(t *testing.T) {
	fn := createVersionExtractor(toold.Linux, toold.Amd64)
	v, ok := fn("apache-maven-3.6.3-bin.tar.gz")
	require.True(t, ok)
	require.Equal(t, toold.ArbitraryVersion{3, 6, 3}, toold.ParseArbitraryVersion(v))
}
