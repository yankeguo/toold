package toold

import (
	"errors"
	"regexp"

	"github.com/yankeguo/numver"
)

// CreateRegexpVersionExtractor creates a version extractor from a regexp layout
func CreateRegexpVersionExtractor(layout string) (fn numver.VersionExtractor, err error) {
	var re *regexp.Regexp
	if re, err = regexp.Compile(layout); err != nil {
		return
	}

	idx := re.SubexpIndex("version")
	if idx == -1 {
		err = errors.New("regexp must have a named capture group: version")
	}

	fn = func(src string) (ver string, ok bool) {
		matches := re.FindStringSubmatch(src)
		if len(matches) == 0 {
			return
		}
		ver = matches[idx]
		ok = true
		return
	}
	return
}
