package toold

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRegexpVersionExtractor(t *testing.T) {
	_, err := CreateRegexpVersionExtractor(`^apache-maven-(?P<ersion>.+).tar.gz$`)
	require.Error(t, err)

	_, err = CreateRegexpVersionExtractor(`^(`)
	require.Error(t, err)

	fn, err := CreateRegexpVersionExtractor(`^apache-maven-(?P<version>.+).tar.gz$`)
	require.NoError(t, err)

	ver, ok := fn("apache-maven-3.6.3.tar.gz")
	require.True(t, ok)
	require.Equal(t, "3.6.3", ver)

	ver, ok = fn("apache-maven-3.8.6.3.tar.gz")
	require.True(t, ok)
	require.Equal(t, "3.8.6.3", ver)

	_, ok = fn("apache-not-maven-3.8.6.3.tar.gz")
	require.False(t, ok)
}
