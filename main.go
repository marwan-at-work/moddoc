package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	// embedded files
	_ "marwan.io/moddoc/statik"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/statik/fs"
	"marwan.io/moddoc/fetch"
	"marwan.io/moddoc/gocopy/semver"
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
		"latestVer":  latestVer,
	}).ParseGlob("frontend/templates/*.html"))
}

func parse() http.FileSystem {
	root := template.New("root").Funcs(template.FuncMap{
		"toLower":    strings.ToLower,
		"subOne":     subOne,
		"getVerLink": getVerLink,
		"json":       getJSON,
		"latestVer":  latestVer,
	})
	dist, err := fs.New()
	must(err)
	hf, err := dist.Open("/templates")
	must(err)
	dir, err := hf.Readdir(-1)
	must(err)
	for _, fi := range dir {
		f, err := dist.Open("/templates/" + fi.Name())
		must(err)
		defer f.Close()
		bts, err := ioutil.ReadAll(f)
		must(err)
		root, err = root.New(fi.Name()).Parse(string(bts))
		must(err)
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
		mods, err := getCatalogModules(r.Context())
		if err != nil {
			fmt.Printf("Error while retrieving catalog from proxy: [%s]\nFallback to public index", err)
			mods, _ = index(r.Context())
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

func getCatalogModules(ctx context.Context) ([]*moduleIndex, error) {
	url := strings.TrimSuffix(config.GoProxyURL, "/") + "/catalog"
	resp, err := fetch.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %v", resp.StatusCode)
	}
	var lr listResp
	err = json.NewDecoder(resp.Body).Decode(&lr)
	if err != nil {
		return nil, err
	}
	mp := map[string][]string{}
	for _, m := range lr.Modules {
		mp[m.Module] = append(mp[m.Module], m.Version)
	}
	mods := []*moduleIndex{}
	for mod, vers := range mp {
		mods = append(mods, &moduleIndex{
			mod,
			vers,
			latestVer(vers),
		})
	}
	return mods, nil
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

func latestVer(vers []string) string {
	sortVersions(vers)
	if len(vers) == 0 {
		return "latest"
	}
	return vers[0]
}

func sortVersions(list []string) {
	sort.Slice(list, func(i, j int) bool {
		cmp := semver.Compare(list[i], list[j])
		if cmp != 0 {
			return cmp < 0
		}
		return list[i] < list[j]
	})
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
