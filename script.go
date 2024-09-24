package toold

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
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

func (sb *ScriptBuilder) AddEcho(s string) {
	s = strings.TrimSpace(s) + "\n"
	sb.b.WriteString("echo -n ")
	sb.b.WriteString(strconv.Quote(base64.StdEncoding.EncodeToString([]byte(s))))
	sb.b.WriteString(" | base64 -d > /dev/stderr\n")
}
