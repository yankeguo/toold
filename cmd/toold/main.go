package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yankeguo/rg"
	"github.com/yankeguo/toold"
	"github.com/yankeguo/toold/adapters/cnpm"
	"github.com/yankeguo/toold/adapters/jdk"
	"github.com/yankeguo/toold/adapters/maven"
	"github.com/yankeguo/toold/adapters/node"
	"github.com/yankeguo/toold/adapters/yarn"
)

var (
	AppVersion = "unknown"
)

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

	log.Println("toold", AppVersion)

	opts := rg.Must(toold.LoadOptions())

	stor := rg.Must(toold.NewStorage(opts))

	hand := toold.NewApp(stor, map[string]toold.Adapter{
		"node":  &node.Adapter{},
		"yarn":  &yarn.Adapter{},
		"cnpm":  &cnpm.Adapter{},
		"jdk":   &jdk.Adapter{},
		"maven": &maven.Adapter{},
	})

	s := &http.Server{
		Addr:    opts.Listen,
		Handler: hand,
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
