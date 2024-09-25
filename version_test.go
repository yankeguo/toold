package toold

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArbitraryVersion(t *testing.T) {
	var cases = []struct {
		v string
		e ArbitraryVersion
	}{
		{"", ArbitraryVersion{}},
		{"1.2.3", ArbitraryVersion{1, 2, 3}},
		{" as 234 35_22  ", ArbitraryVersion{234, 35, 22}},
		{"1.0.3", ArbitraryVersion{1, 0, 3}},
		{"1.0002.3", ArbitraryVersion{1, 2, 3}},
	}
	for _, c := range cases {
		require.Equal(t, c.e, ParseArbitraryVersion(c.v))
	}
}

func TestArbitraryVersionMatch(t *testing.T) {
	var cases = []struct {
		v ArbitraryVersion
		c ArbitraryVersion
		e bool
	}{
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 3}, true},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2}, true},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 4}, false},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 3, 4}, false},
	}
	for _, c := range cases {
		require.Equal(t, c.e, c.v.Match(c.c))
	}
}

func TestArbitraryVersionString(t *testing.T) {
	var cases = []struct {
		v ArbitraryVersion
		e string
	}{
		{ArbitraryVersion{1, 2, 3}, "1.2.3"},
		{ArbitraryVersion{234, 35, 22}, "234.35.22"},
		{ArbitraryVersion{1, 0, 30}, "1.0.30"},
	}
	for _, c := range cases {
		require.Equal(t, c.e, c.v.String())
	}
}

func TestArbitraryVersionLessThan(t *testing.T) {
	var cases = []struct {
		v ArbitraryVersion
		o ArbitraryVersion
		e int
	}{
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 3}, 0},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 4}, -1},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 2}, 1},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2}, 1},
		{ArbitraryVersion{1, 2, 3}, ArbitraryVersion{1, 2, 3, 4}, -1},
	}
	for _, c := range cases {
		require.Equal(t, c.e, c.v.Compare(c.o))
	}
}

func TestFindBestVersionedFile(t *testing.T) {
	bf, bv, err := FindBestVersionedFile(
		FindBestVersionedFileOptions{
			Files: []string{
				"deno-v1.2.3-linux-x64.tar.gz",
				"node-v1.2.3-linux-x64.tar.gz",
				"node-v1.2.4-windows-x64.tar.gz",
				"node-v1.3.3-linux-x64.tar.gz",
				"node-v1.3.4-windows-x64.tar.gz",
			},
			VersionExtractor: func(s string) (v string, ok bool) {
				const (
					prefix = "node-v"
					suffix = "-linux-x64.tar.gz"
				)
				ok = strings.HasPrefix(s, prefix) && strings.HasSuffix(s, suffix)
				v = strings.TrimPrefix(strings.TrimSuffix(s, suffix), prefix)
				return
			},
			VersionConstraint: "1.3",
		},
	)
	require.NoError(t, err)
	require.Equal(t, "node-v1.3.3-linux-x64.tar.gz", bf)
	require.Equal(t, ArbitraryVersion{1, 3, 3}, bv)
}
