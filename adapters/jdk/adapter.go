package jdk

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
		toold.Darwin: "mac",
		toold.Linux:  "linux",
	}
	mArch = map[string]string{
		toold.Amd64: "x64",
		toold.Arm64: "aarch64",
	}
)

func CreateVersionExtractor(opts toold.AdapterOptions) numver.VersionExtractor {
	return rg.Must(toold.CreateRegexpVersionExtractor(
		fmt.Sprintf(
			`jdk_%s_%s_hotspot_(?<version>.+)\.tar\.gz$`,
			toold.ResolvePlatform(opts.Arch, mArch),
			toold.ResolvePlatform(opts.OS, mOS),
		),
	))
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "jdk"))

	file, version, found := numver.Search(numver.SearchOptions{
		Items:      files,
		Constraint: opts.Version,
		Extractor:  CreateVersionExtractor(opts),
	})

	if !found {
		err = errors.New("jdk version not found for: " + opts.Version)
		return
	}

	envJavaHome := ""
	envPrependPath := "/bin"

	if opts.OS == toold.Darwin {
		envJavaHome = "/Contents/Home"
		envPrependPath = "/Contents/Home/bin"
	}

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             rg.Must(opts.Storage.CreateSignedURL(ctx, "jdk/"+file, time.Minute*10)),
		Dir:             "jdk-" + version.String(),
		StripComponents: 1,
		Env: map[string]string{
			"JAVA_HOME": envJavaHome,
		},
		EnvPrependPath: []string{envPrependPath},
	})
	return
}
