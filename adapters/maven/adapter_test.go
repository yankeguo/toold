package maven

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/numver"
	"github.com/yankeguo/toold"
)

func TestCreateVersionExtractor(t *testing.T) {
	fn := CreateVersionExtractor(toold.AdapterOptions{})
	v, ok := fn("apache-maven-3.6.3-bin.tar.gz")
	require.True(t, ok)
	require.Equal(t, numver.Version{3, 6, 3}, numver.Parse(v))
}
