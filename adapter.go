package toold

import "context"

type AdapterOptions struct {
	Storage *Storage
	Script  *ScriptBuilder
	OS      string
	Arch    string
	Name    string
	Version string
	Force   bool
}

type Adapter interface {
	Build(ctx context.Context, opts AdapterOptions) (err error)
}
