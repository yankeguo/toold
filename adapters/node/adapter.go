package node

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yankeguo/numver"
	"github.com/yankeguo/rg"
	"github.com/yankeguo/toold"
)

var (
	mOS = map[string]string{
		toold.Darwin: "darwin",
		toold.Linux:  "linux",
	}
	mArch = map[string]string{
		toold.Amd64: "x64",
		toold.Arm64: "arm64",
	}
)

func CreateVersionExtractor(opts toold.AdapterOptions) numver.VersionExtractor {
	return rg.Must(toold.CreateRegexpVersionExtractor(
		fmt.Sprintf(
			`^node-v(?<version>.+)-%s-%s\.tar\.gz$`,
			toold.ResolvePlatform(opts.OS, mOS),
			toold.ResolvePlatform(opts.Arch, mArch),
		),
	))
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "node"))

	file, version, found := numver.Search(numver.SearchOptions{
		Items:      files,
		Constraint: opts.Version,
		Extractor:  CreateVersionExtractor(opts),
	})

	if !found {
		err = errors.New("node version not found for: " + opts.Version)
		return
	}

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             rg.Must(opts.Storage.CreateSignedURL(ctx, "node/"+file, time.Minute*10)),
		Dir:             "node-" + version.String(),
		StripComponents: 1,
		EnvPrependPath:  []string{"/bin"},
	})
	return
}
