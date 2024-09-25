package maven

import (
	"context"
	"strings"
	"time"

	"github.com/yankeguo/rg"
	"github.com/yankeguo/toold"
)

func createVersionExtractor(os string, arch string) toold.VersionExtractor {
	return func(src string) (ver string, ok bool) {
		// check tar.gz
		if !strings.HasSuffix(src, ".tar.gz") {
			return
		}
		if !strings.Contains(src, "maven") {
			return
		}
		ver = src
		ok = true
		return
	}
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "maven"))

	file, version := rg.Must2(toold.FindBestVersionedFile(toold.FindBestVersionedFileOptions{
		Files:             files,
		VersionExtractor:  createVersionExtractor(opts.OS, opts.Arch),
		VersionConstraint: opts.Version,
	}))

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             rg.Must(opts.Storage.CreateSignedURL(ctx, "maven/"+file, time.Minute*10)),
		Dir:             "maven-" + version.String(),
		StripComponents: 1,
		EnvPrependPath:  []string{"/bin"},
	})
	return
}
