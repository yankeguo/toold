package maven

import (
	"context"
	"errors"
	"time"

	"github.com/yankeguo/numver"
	"github.com/yankeguo/rg"
	"github.com/yankeguo/toold"
)

func CreateVersionExtractor(opts toold.AdapterOptions) numver.VersionExtractor {
	return rg.Must(toold.CreateRegexpVersionExtractor(
		`maven-(?<version>.+)-bin\.tar\.gz$`,
	))
}

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	defer rg.Guard(&err)

	files := rg.Must(opts.Storage.ListFiles(ctx, "maven"))

	file, version, found := numver.Search(numver.SearchOptions{
		Items:      files,
		Constraint: opts.Version,
		Extractor:  CreateVersionExtractor(opts),
	})

	if !found {
		err = errors.New("maven version not found for: " + opts.Version)
		return
	}

	opts.Out.AddDownloadAndExtract(toold.ScriptDownloadAndExtractOptions{
		URL:             rg.Must(opts.Storage.CreateSignedURL(ctx, "maven/"+file, time.Minute*10)),
		Dir:             "maven-" + version.String(),
		StripComponents: 1,
		EnvPrependPath:  []string{"/bin"},
	})
	return
}
