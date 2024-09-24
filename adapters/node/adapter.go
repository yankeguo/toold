package node

import (
	"context"
	"errors"
	"time"

	"github.com/yankeguo/rg"
	"github.com/yankeguo/toold"
)

func createFilePattern(os string, arch string) (prefix string, suffix string, ok bool) {
	_os, _arch := os, arch
	switch os {
	case toold.Darwin:
	case toold.Linux:
	default:
		return
	}
	switch arch {
	case toold.Amd64:
		_arch = "x64"
	case toold.Arm64:
	default:
		return
	}
	return "node-v", "-" + _os + "-" + _arch + ".tar.gz", true
}

type Adapter struct{}

func (a Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	pfx, sfx, ok := createFilePattern(opts.OS, opts.Arch)
	if !ok {
		err = errors.New("unsupported os/arch: " + opts.OS + "/" + opts.Arch)
		return
	}

	files := rg.Must(opts.Storage.ListFiles(ctx, "node"))

	file, version := rg.Must2(toold.FindBestVersionedFile(toold.FindBestVersionedFileOptions{
		Files:   files,
		Prefix:  pfx,
		Suffix:  sfx,
		Version: opts.Version,
	}))

	localDir := "node-" + version.String()

	link := rg.Must(opts.Storage.CreateSignedURL(ctx, "node/"+file, time.Minute*10))

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             link,
		Dir:             localDir,
		StripComponents: 1,
		PrependPath:     "bin",
	})
	return
}
