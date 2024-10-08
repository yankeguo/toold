package cnpm

import (
	"context"

	"github.com/yankeguo/toold"
)

type Adapter struct{}

func (a *Adapter) Build(ctx context.Context, opts toold.AdapterOptions) (err error) {
	opts.Script.AddScriptGlobalNodePackageOptions(
		toold.ScriptGlobalNodePackageOptions{
			Command:  "cnpm",
			Package:  "cnpm",
			Registry: "https://registry.npmmirror.com",
			Version:  opts.Version,
			Force:    opts.Force,
		},
	)
	return
}
