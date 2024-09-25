package toold

func ResolvePlatform(s string, m map[string]string) string {
	if v, ok := m[s]; ok {
		return v
	}
	return "unsupported"
}
