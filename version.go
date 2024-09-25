package toold

import (
	"errors"
	"strconv"
	"strings"
)

type ArbitraryVersion []int

func ParseArbitraryVersion(v string) ArbitraryVersion {
	sv := ArbitraryVersion{}

	var chunk string
	for _, c := range v {
		if c <= '9' && c >= '0' {
			chunk += string(c)
		} else {
			if chunk != "" {
				n, _ := strconv.Atoi(chunk)
				sv = append(sv, n)
				chunk = ""
			}
		}
	}

	if chunk != "" {
		n, _ := strconv.Atoi(chunk)
		sv = append(sv, n)
	}

	return sv
}

func (sv ArbitraryVersion) String() string {
	var parts []string
	for _, n := range sv {
		parts = append(parts, strconv.Itoa(n))
	}
	return strings.Join(parts, ".")
}

func (sv ArbitraryVersion) Match(constraint ArbitraryVersion) bool {
	if len(constraint) == 0 {
		return true
	}
	if len(constraint) > len(sv) {
		return false
	}
	for i, n := range constraint {
		if sv[i] != n {
			return false
		}
	}
	return true
}

func (sv ArbitraryVersion) Compare(other ArbitraryVersion) int {
	for i := 0; i < len(sv) && i < len(other); i++ {
		if sv[i] < other[i] {
			return -1
		}
		if sv[i] > other[i] {
			return 1
		}
	}
	switch {
	case len(sv) < len(other):
		return -1
	case len(sv) > len(other):
		return 1
	default:
		return 0
	}
}

type VersionExtractor func(src string) (ver string, ok bool)

type FindBestVersionedFileOptions struct {
	Files             []string
	VersionExtractor  VersionExtractor
	VersionConstraint string
}

func FindBestVersionedFile(opts FindBestVersionedFileOptions) (bestFile string, bestVersion ArbitraryVersion, err error) {
	constraint := ParseArbitraryVersion(opts.VersionConstraint)

	type eligibleItem struct {
		file    string
		version ArbitraryVersion
	}
	var items []eligibleItem

	for _, file := range opts.Files {
		// version extractor
		v, ok := opts.VersionExtractor(file)
		if !ok {
			continue
		}
		// av constraint
		av := ParseArbitraryVersion(v)
		if av.Match(constraint) {
			items = append(items, eligibleItem{
				file:    file,
				version: av,
			})
		}
	}

	if len(items) == 0 {
		err = errors.New("no matching file found for version: " + opts.VersionConstraint)
		return
	}

	for _, item := range items {
		if bestVersion == nil || item.version.Compare(bestVersion) > 0 {
			bestFile = item.file
			bestVersion = item.version
		}
	}

	return
}
