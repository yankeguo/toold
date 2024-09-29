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

func (sb *ScriptBuilder) Add(s string) {
	sb.b.WriteString(s)
}

func (sb *ScriptBuilder) AddTemplate(layout string, data map[string]any) {
	if !strings.HasSuffix(layout, "\n") {
		layout += "\n"
	}
	tpl := template.Must(template.New("script").Parse(layout))
	buf := &bytes.Buffer{}
	rg.Must0(tpl.Execute(buf, data))
	sb.b.WriteString(buf.String())
}

func (sb *ScriptBuilder) AddWarning(s string) {
	items := strings.Split(s, "\n")
	for i := range items {
		items[i] = "toold: " + items[i]
	}

	sb.AddTemplate(
		`echo "{{.content}}" | base64 -d 1>&2`,
		map[string]any{
			"content": base64.StdEncoding.EncodeToString([]byte(strings.Join(items, "\n"))),
		},
	)
}

type ScriptGlobalNodePackageOptions struct {
	Command  string
	Package  string
	Registry string
	Version  string
	Force    bool
}

func (sb *ScriptBuilder) AddScriptGlobalNodePackageOptions(opts ScriptGlobalNodePackageOptions) {
	sb.AddTemplate(`
{{if .force}}
{{else}}
if command -v "{{.command}}" > /dev/null; then
    echo "toold: {{.command}} is already installed" 1>&2
else
{{end}}
    echo "toold: installing {{.command}}" 1>&2
    npm install -g "{{.package}}{{if .version}}@{{.version}}{{end}}"{{if .registry}} --registry="{{.registry}}"{{end}} 1>&2
{{if .force}}
{{else}}
fi
{{end}}
`, map[string]any{
		"command":  opts.Command,
		"package":  opts.Package,
		"registry": opts.Registry,
		"version":  opts.Version,
		"force":    opts.Force,
	})
}

type ScriptDownloadAndExtractOptions struct {
	URL             string
	Dir             string
	StripComponents int
	EnvPrependPath  []string
	Env             map[string]string
}

func (sb *ScriptBuilder) AddDownloadAndExtract(opts ScriptDownloadAndExtractOptions) {
	sb.AddTemplate(`
TOOLD_HOME="${TOOLD_HOME}"
if [ -z "${TOOLD_HOME}" ]; then
    echo "toold: \$TOOLD_HOME is not set, using ~/.toold" 1>&2
    TOOLD_HOME="${HOME}/.toold"
fi
mkdir -p "${TOOLD_HOME}" 1>&2

TOOL_DIR="${TOOLD_HOME}/{{.dir}}"
PID_FILE="${TOOL_DIR}.incomplete"

while true; do
    if [ -f "${PID_FILE}" ]; then
        PID=$(cat "${PID_FILE}")
        if [ -n "${PID}" ] && kill -0 "${PID}" &> /dev/null; then
            echo "toold: waiting for another process to finish" 1>&2
            sleep 5
        else
            echo "toold: found incomplete dir {{.dir}}, cleaning up" 1>&2
            rm -rf "${PID_FILE}" "${TOOL_DIR}" 1>&2
            break
        fi
    else
        break
    fi
done

if [ ! -d "${TOOL_DIR}" ]; then
    echo -n $$ > "${PID_FILE}"
    mkdir -p "${TOOL_DIR}" 1>&2
    echo "toold: downloading {{.dir}}" 1>&2
    curl -sSL "{{.url}}" | tar -xz -C "${TOOL_DIR}" {{if .strip_components}}--strip-components={{.strip_components}}{{end}} 1>&2
    rm -f "${PID_FILE}" 1>&2
fi

{{if .env_prepend_path}}
echo "toold: using {{.dir}}" 1>&2
{{range .env_prepend_path}}
export PATH="${TOOL_DIR}{{.}}:$PATH"
echo "export PATH=\"${TOOL_DIR}{{.}}:\$PATH\""
{{end}}
{{end}}

{{if .env}}
{{range $key, $value := .env}}
export {{$key}}="${TOOL_DIR}{{$value}}"
echo "export {{$key}}=\"${TOOL_DIR}{{$value}}\""
{{end}}
{{end}}
`,
		map[string]any{
			"dir":              opts.Dir,
			"strip_components": opts.StripComponents,
			"url":              opts.URL,
			"env_prepend_path": opts.EnvPrependPath,
			"env":              opts.Env,
		},
	)
}
