package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"marwan.io/moddoc/fetch"
)

type listResp struct {
	Modules []module `json:"modules"`
	Next    string   `json:"next"`
}

type module struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

type moduleIndex struct {
	Module   string   `json:"module"`
	Versions []string `json:"versions"`
	Latest   string   `json:"latest"`
}

func catalog(w http.ResponseWriter, r *http.Request) {
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
		mods = append(mods, &moduleIndex{
			mod,
			vers,
			latestVer(vers),
		})
	}
	json.NewEncoder(w).Encode(mods)
}
