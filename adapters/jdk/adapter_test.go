package jdk

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yankeguo/numver"
	"github.com/yankeguo/toold"
)

func TestCreateVersionExtractor(t *testing.T) {
	fn := CreateVersionExtractor(toold.AdapterOptions{OS: toold.Linux, Arch: toold.Amd64})
	v, ok := fn("OpenJDK11U-jdk_x64_linux_hotspot_11.0.24_8.tar.gz")
	require.True(t, ok)
	require.Equal(t, numver.Version{11, 0, 24, 8}, numver.Parse(v))
}
