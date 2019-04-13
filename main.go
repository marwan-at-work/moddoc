package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	// embedded files
	_ "marwan.io/moddoc/statik"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/statik/fs"
	"marwan.io/moddoc/fetch"
	"marwan.io/moddoc/proxy"
)

//go:generate statik -src=frontend
var config struct {
	GoProxyURL string `envconfig:"GOPROXY" required:"true"`
	Port       string `envconfig:"PORT" default:"3001"`
	ENV        string `envconfig:"MODDOC_ENV"`
}

func init() {
	envconfig.MustProcess("", &config)
}

var tt *template.Template

func parseDev() {
	tt = template.Must(template.New("root").Funcs(template.FuncMap{
		"toLower":    strings.ToLower,
		"subOne":     subOne,
		"getVerLink": getVerLink,
		"json":       getJSON,
	}).ParseGlob("frontend/templates/*.html"))
}

func parse() http.FileSystem {
	root := template.New("root").Funcs(template.FuncMap{
		"toLower":    strings.ToLower,
		"subOne":     subOne,
		"getVerLink": getVerLink,
		"json":       getJSON,
	})
	dist, err := fs.New()
	if err != nil {
		panic(err)
	}
	hf, err := dist.Open("/templates")
	if err != nil {
		panic(err)
	}
	dir, err := hf.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, fi := range dir {
		f, err := dist.Open("/templates/" + fi.Name())
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bts, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		root, err = root.New(fi.Name()).Parse(string(bts))
		if err != nil {
			panic(err)
		}
	}
	tt = root
	return dist

}

func main() {
	r := mux.NewRouter()
	srv := proxy.NewService(config.GoProxyURL)
	dist := parse()
	r.HandleFunc("/", home(dist))
	r.HandleFunc("/catalog", catalog)
	r.Handle(docPath, getDoc(srv))
	if config.ENV == "DEV" {
		parseDev()
		r.PathPrefix("/public/").Handler(http.FileServer(http.Dir("frontend")))
	} else {
		r.PathPrefix("/public/").Handler(http.FileServer(dist))
	}

	fmt.Println("listening on port :" + config.Port)
	http.ListenAndServe(":"+config.Port, r)
}

func home(fs http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimSuffix(config.GoProxyURL, "/") + "/catalog"
		resp, err := fetch.Fetch(r.Context(), url)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			http.Error(w, "unexpected status: "+resp.Status, resp.StatusCode)
			return
		}
		var lr listResp
		err = json.NewDecoder(resp.Body).Decode(&lr)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		mp := map[string][]string{}
		for _, m := range lr.Modules {
			mp[m.Module] = append(mp[m.Module], m.Version)
		}
		mods := []*moduleIndex{}
		for mod, vers := range mp {
			mods = append(mods, &moduleIndex{mod, vers})
		}
		err = tt.Lookup("index.html").Execute(w, map[string]interface{}{
			"index": true,
			"data":  mods,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func subOne(i int) int {
	return i - 1
}

func getVerLink(importPath, version string) string {
	return filepath.Join("/", importPath, "@v", version)
}

func getJSON(i interface{}) string {
	bts, _ := json.Marshal(i)
	return string(bts)
}
