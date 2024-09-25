package node

import (
	"context"
	"strings"
	"time"

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

func createVersionExtractor(os string, arch string) toold.VersionExtractor {
	platform := "-" + toold.ResolvePlatform(os, mOS) + "-" + toold.ResolvePlatform(arch, mArch)

	return func(src string) (ver string, ok bool) {
		// check tar.gz
		if !strings.HasSuffix(src, ".tar.gz") {
			return
		}
		// check platform
		if !strings.Contains(src, platform) {
			return
		}
		// extract before platform
		ver = src[0:strings.Index(src, platform)]
		ok = true
		return
	}
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "node"))

	file, version := rg.Must2(toold.FindBestVersionedFile(toold.FindBestVersionedFileOptions{
		Files:             files,
		VersionExtractor:  createVersionExtractor(opts.OS, opts.Arch),
		VersionConstraint: opts.Version,
	}))

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             rg.Must(opts.Storage.CreateSignedURL(ctx, "node/"+file, time.Minute*10)),
		Dir:             "node-" + version.String(),
		StripComponents: 1,
		EnvPrependPath:  "bin",
	})
	return
}
