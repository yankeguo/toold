package yarn

import (
	"context"

	"github.com/yankeguo/toold"
)

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	opts.Script.AddScriptGlobalNodePackageOptions(
		toold.ScriptGlobalNodePackageOptions{
			Command: "yarn",
			Package: "yarn",
			Version: opts.Version,
			Force:   opts.Force,
		},
	)
	return
}
