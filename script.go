package toold

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/yankeguo/rg"
)

type ScriptBuilder struct {
	b *strings.Builder
}

func NewScriptBuilder() *ScriptBuilder {
	return &ScriptBuilder{b: &strings.Builder{}}
}

func (sb *ScriptBuilder) Reset() {
	sb.b.Reset()
}

func (sb *ScriptBuilder) Concat(sub *ScriptBuilder) {
	sb.b.WriteString(sub.b.String())
}

func (sb *ScriptBuilder) WriteTo(rw http.ResponseWriter) {
	buf := []byte(sb.b.String())
	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	rw.Write(buf)
}

func (sb *ScriptBuilder) add(layout string, data map[string]any) {
	if !strings.HasSuffix(layout, "\n") {
		layout += "\n"
	}
	tpl := template.Must(template.New("script").Parse(layout))
	buf := &bytes.Buffer{}
	rg.Must0(tpl.Execute(buf, data))
	sb.b.WriteString(buf.String())
}

func (sb *ScriptBuilder) AddEcho(s string) {
	s = strings.TrimSpace(s) + "\n"

	sb.add(
		`echo -n "{{.content}}" | base64 -d > /dev/stderr`,
		map[string]any{
			"content": base64.StdEncoding.EncodeToString([]byte(s)),
		},
	)
}

type ScriptDownloadAndExtractOptions struct {
	URL             string
	Dir             string
	StripComponents int
	PrependPath     string
}

func (sb *ScriptBuilder) AddDownloadAndExtract(opts ScriptDownloadAndExtractOptions) {
	sb.add(`
TOOLD_ROOT="${TOOLD_ROOT}"
if [ -z "${TOOLD_ROOT}" ]; then
  echo "TOOLD_ROOT is not set, default to ~/.toold" > /dev/stderr
  TOOLD_ROOT="${HOME}/.toold"
fi

mkdir -p "${TOOLD_ROOT}"

if [ -f "${TOOLD_HOME}/{{.dir}}.incomplete" ]; then
  echo "found incomplete dir {{.dir}}, cleaning up" > /dev/stderr
  rm -rf "${TOOLD_ROOT}/{{.dir}}"
fi

touch "${TOOLD_ROOT}/{{.dir}}.incomplete"

mkdir -p "${TOOLD_ROOT}/{{.dir}}"

curl -sSL "{{.url}}" | tar -xz -C "${TOOLD_ROOT}/{{.dir}}" {{if .strip_components}}--strip-components={{.strip_components}}{{end}}

rm -f "${TOOLD_ROOT}/{{.dir}}.incomplete"

echo "{{.dir}} completed" > /dev/stderr

{{if .prepend_path}}
echo "export PATH=\"${TOOLD_ROOT}/{{.dir}}/{{.prepend_path}}:\$PATH\""
{{end}}
`,
		map[string]any{
			"dir":              opts.Dir,
			"strip_components": opts.StripComponents,
			"url":              opts.URL,
			"prepend_path":     opts.PrependPath,
		},
	)
}
