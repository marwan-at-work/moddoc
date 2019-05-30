package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"marwan.io/moddoc/fetch"
	gomodule "marwan.io/moddoc/gocopy/module"
)

func getModule(w http.ResponseWriter, r *http.Request) {
	mod := strings.TrimPrefix(r.URL.Path, "/")
	mod, err := gomodule.EncodePath(strings.TrimSuffix(mod, "/"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	url := config.GoProxyURL + "/" + mod + "/@v/list"
	resp, err := fetch.Fetch(r.Context(), url)
	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching list: %v", err), 500)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		http.NotFound(w, r)
		return
	}
	scnr := bufio.NewScanner(resp.Body)
	vers := []string{}
	for scnr.Scan() {
		vers = append(vers, scnr.Text())
	}
	ver := latestVer(vers)
	if ver == "latest" {
		url := config.GoProxyURL + "/" + mod + "/@latest"
		ver = getLatest(url)
	}
	http.Redirect(w, r, "/"+mod+"/@v/"+ver, http.StatusMovedPermanently)
}

func getLatest(url string) string {
	resp, err := fetch.Fetch(context.Background(), url)
	if err != nil {
		return "latest"
	}
	defer resp.Body.Close()
	var mod struct {
		Version string
	}
	json.NewDecoder(resp.Body).Decode(&mod)
	if mod.Version == "" {
		mod.Version = "latest"
	}
	return mod.Version
}
