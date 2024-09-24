package toold

import (
	"errors"
	"flag"
	"os"

	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"
)

const (
	BackendCOS = "cos"
)

type Options struct {
	Verbose bool   `yaml:"verbose"`
	Listen  string `yaml:"listen"`
	Backend string `yaml:"backend"`
	COS     struct {
		BucketURL string `yaml:"bucket_url"`
		SecretID  string `yaml:"secret_id"`
		SecretKey string `yaml:"secret_key"`
	} `yaml:"cos"`
}

func LoadOptions() (opts Options, err error) {
	defer rg.Guard(&err)

	// load options from config file
	var conf string
	flag.StringVar(&conf, "conf", "toold.yaml", "config file")
	flag.Parse()

	// unmarshal options
	rg.Must0(yaml.Unmarshal(rg.Must(os.ReadFile(conf)), &opts))

	if opts.Listen == "" {
		opts.Listen = ":8080"
	}

	if opts.Backend == "" {
		opts.Backend = BackendCOS
	}

	if opts.Backend != BackendCOS {
		err = errors.New("unknown backend: " + opts.Backend)
		return
	}

	if opts.COS.BucketURL == "" {
		err = errors.New("missing cos.bucket_url")
		return
	}

	if opts.COS.SecretID == "" {
		err = errors.New("missing cos.secret_id")
		return
	}

	if opts.COS.SecretKey == "" {
		err = errors.New("missing cos.secret_key")
		return
	}

	return
}
