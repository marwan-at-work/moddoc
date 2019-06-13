package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
	config.GoProxyURL = parseProxyURL(config.GoProxyURL)
}

var tt *template.Template

func parseDev() {
	tt = template.Must(template.New("root").Funcs(template.FuncMap{
		"toLower":        strings.ToLower,
		"subOne":         subOne,
		"getVerLink":     getVerLink,
		"json":           getJSON,
		"latestVer":      latestVer,
		"methodReceiver": methodReceiver,
	}).ParseGlob("frontend/templates/*.html"))
}

func parse() http.FileSystem {
	root := template.New("root").Funcs(template.FuncMap{
		"toLower":        strings.ToLower,
		"subOne":         subOne,
		"getVerLink":     getVerLink,
		"json":           getJSON,
		"latestVer":      latestVer,
		"methodReceiver": methodReceiver,
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

func parseProxyURL(s string) string {
	proxyURL := strings.Split(s, ",")[0]
	switch proxyURL {
	case "":
		log.Fatal("GOPROXY's first argument must not be empty")
	case "direct":
		log.Fatal("cannot use 'direct' as a GOPROXY destination")
	case "off":
		log.Fatal("cannot use 'off' as a GOPROXY destination")
	}
	return proxyURL
}

func main() {
	r := mux.NewRouter()
	srv := proxy.NewService(config.GoProxyURL)
	dist := parse()
	r.Handle("/", home(dist))
	r.Handle(docPath, getDoc(srv))
	r.HandleFunc("/catalog", catalog)
	if config.ENV == "DEV" {
		parseDev()
		r.PathPrefix("/public/").Handler(http.FileServer(http.Dir("frontend")))
	} else {
		r.PathPrefix("/public/").Handler(http.FileServer(dist))
	}
	r.NotFoundHandler = http.HandlerFunc(getModule)

	fmt.Println("listening on port :" + config.Port)
	http.ListenAndServe(":"+config.Port, r)
}

func home(fs http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mods, err := getCatalogModules(r.Context())
		if err != nil {
			fmt.Printf("Error while retrieving catalog from proxy: [%s]\nFallback to public index\n", err)
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
	sort.Slice(list, func(i, j int) bool { return semver.Compare(list[i], list[j]) > 0 })
}

func methodReceiver(receiver string) string {
	if receiver == "" {
		return ""
	}

	return "(" + receiver + ")"
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
