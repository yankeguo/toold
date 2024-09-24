package toold

import (
	"log"
	"net/http"
)

type App struct {
	Storage *Storage
}

func (h *App) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m := ParseManifest(r.URL.Path)
	log.Println("manifest:", m.String())
	rw.Write([]byte(""))
}
