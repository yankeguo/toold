package toold

import "context"

type AdapterOptions struct {
	Storage *Storage
	Out     *ScriptBuilder
	OS      string
	Arch    string
	Name    string
	Version string
}

type Adapter interface {
	Build(ctx context.Context, opts AdapterOptions) (err error)
}
