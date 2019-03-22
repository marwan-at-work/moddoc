package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	// embedded files
	_ "marwan.io/moddoc/statik"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/statik/fs"
	"marwan.io/moddoc/proxy"
)

var config struct {
	GoProxyURL string `envconfig:"GOPROXY" required:"true"`
	Port       string `envconfig:"PORT" default:"3001"`
}

func init() {
	envconfig.MustProcess("", &config)
}

func main() {
	r := mux.NewRouter()
	srv := proxy.NewService(config.GoProxyURL)
	dist, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r.HandleFunc("/", home(dist))
	r.HandleFunc("/catalog", catalog)
	r.Handle(docPath, getDoc(srv))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer(dist)))
	r.NotFoundHandler = http.HandlerFunc(home(dist))

	fmt.Println("listening on port :" + config.Port)
	http.ListenAndServe(":"+config.Port, r)
}

func home(fs http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open("/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		defer f.Close()
		io.Copy(w, f)
	}
}
