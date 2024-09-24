package toold

import "strings"

var (
	osDarwin  = "darwin"
	osWindows = "windows"
	osLinux   = "linux"

	archAmd64 = "amd64"
	archArm64 = "arm64"

	supportedOS   = []string{osLinux, osDarwin, osWindows}
	supportedArch = []string{archAmd64, archArm64}
)

type ManifestTool struct {
	Name    string
	Version string
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
		for _, os := range supportedOS {
			if c == os {
				m.OS = os
				continue outerLoop
			}
		}
		for _, arch := range supportedArch {
			if c == arch {
				m.Arch = arch
				continue outerLoop
			}
		}
		splits := strings.Split(c, "@")
		if len(splits) != 2 {
			continue
		}
		key, ver := strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1])
		if key == "" || ver == "" {
			continue
		}
		m.Tools = append(m.Tools, ManifestTool{
			Name:    key,
			Version: ver,
		})
	}

	if m.OS == "" {
		m.OS = supportedOS[0]
	}
	if m.Arch == "" {
		m.Arch = supportedArch[0]
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
