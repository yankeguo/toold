package toold

import (
	"context"
	"net/http"
)

func NewApp(storage *Storage, adapters map[string]Adapter) *App {
	if adapters == nil {
		adapters = make(map[string]Adapter)
	}
	return &App{
		storage:  storage,
		adapters: adapters,
	}
}

type App struct {
	storage  *Storage
	adapters map[string]Adapter
}

func (h *App) build(ctx context.Context, out *ScriptBuilder, m Manifest) {
	for _, tool := range m.Tools {
		adapter, ok := h.adapters[tool.Name]
		if !ok {
			out.AddEcho("adapter not found for [" + tool.Name + "]")
			continue
		}
		sub := NewScriptBuilder()
		if err := adapter.Build(ctx, AdapterOptions{
			Storage: h.storage,
			Out:     sub,
			OS:      m.OS,
			Arch:    m.Arch,
			Name:    tool.Name,
			Version: tool.Version,
		}); err != nil {
			sub.Reset()
			sub.AddEcho("adapter [" + tool.Name + "] failed: " + err.Error())
		}
		out.Concat(sub)
	}
}

func (h *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m := ParseManifest(r.URL.Path)
	out := NewScriptBuilder()
	h.build(r.Context(), out, m)
	out.WriteTo(rw)
}
