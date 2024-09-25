package jdk

import (
	"context"
	"strings"
	"time"

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

func createVersionExtractor(os string, arch string) toold.VersionExtractor {
	platform := "jdk_" + toold.ResolvePlatform(arch, mArch) + "_" + toold.ResolvePlatform(os, mOS)

	return func(src string) (ver string, ok bool) {
		// check tar.gz
		if !strings.HasSuffix(src, ".tar.gz") {
			return
		}
		// check platform
		if !strings.Contains(src, platform) {
			return
		}
		// extract after version
		ver = src[strings.Index(src, platform)+len(platform):]
		ok = true
		return
	}
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "jdk"))

	file, version := rg.Must2(toold.FindBestVersionedFile(toold.FindBestVersionedFileOptions{
		Files:             files,
		VersionExtractor:  createVersionExtractor(opts.OS, opts.Arch),
		VersionConstraint: opts.Version,
	}))

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
