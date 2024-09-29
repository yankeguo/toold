package toold

import "strings"

const (
	Darwin = "darwin"
	Linux  = "linux"

	Amd64 = "amd64"
	Arm64 = "arm64"
)

var (
	SupportedOS   = []string{Linux, Darwin}
	SupportedArch = []string{Amd64, Arm64}
)

type ManifestTool struct {
	Name    string
	Version string
	Force   bool
}

type Manifest struct {
	OS    string
	Arch  string
	Tools []ManifestTool
}

func ParseManifest(u string) (m Manifest) {
outerLoop:
	for _, c := range strings.Split(u, "/") {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		for _, os := range SupportedOS {
			if c == os {
				m.OS = os
				continue outerLoop
			}
		}
		for _, arch := range SupportedArch {
			if c == arch {
				m.Arch = arch
				continue outerLoop
			}
		}
		splits := strings.Split(c, "@")
		if len(splits) == 1 {
			key := strings.TrimSpace(splits[0])
			if key == "" {
				continue
			}
			m.Tools = append(m.Tools, ManifestTool{
				Name: key,
			})
		} else if len(splits) == 2 {
			key, ver := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
			if key == "" || ver == "" {
				continue
			}
			force := strings.HasSuffix(ver, "!")
			ver = strings.TrimSuffix(ver, "!")
			m.Tools = append(m.Tools, ManifestTool{
				Name:    key,
				Version: ver,
				Force:   force,
			})
		}
	}

	if m.OS == "" {
		m.OS = SupportedOS[0]
	}
	if m.Arch == "" {
		m.Arch = SupportedArch[0]
	}
	return
}

func (m Manifest) String() string {
	sb := &strings.Builder{}
	sb.WriteString("{PLATFORM=")
	sb.WriteString(m.OS)
	sb.WriteString("/")
	sb.WriteString(m.Arch)
	sb.WriteString(", TOOLS=")
	for i, t := range m.Tools {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(t.Name)
		sb.WriteString("@")
		sb.WriteString(t.Version)
	}
	sb.WriteString("}")
	return sb.String()
}
