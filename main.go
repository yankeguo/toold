package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Verbose bool   `yaml:"verbose"`
	Listen  string `yaml:"listen"`
}

func loadOptions() (opts Options, err error) {
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

	return
}

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exit with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	opts := rg.Must(loadOptions())

	m := &http.ServeMux{}

	s := &http.Server{
		Addr:    opts.Listen,
		Handler: m,
	}

	chErr := make(chan error, 1)
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("listening on:", opts.Listen)
		chErr <- s.ListenAndServe()
	}()

	select {
	case err = <-chErr:
		return
	case sig := <-chSig:
		log.Println("signal caught:", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	err = s.Shutdown(ctx)
	<-chErr
}
