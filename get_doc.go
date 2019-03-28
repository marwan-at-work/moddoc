package main

import (
	"net/http"

	"github.com/gorilla/mux"
	gomodule "marwan.io/moddoc/gocopy/module"
	"marwan.io/moddoc/proxy"
)

var docPath = "/{module:.+}/@v/{version}"

func getDoc(proxy proxy.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mod := mux.Vars(r)["module"]
		ver := mux.Vars(r)["version"]
		mod, err := gomodule.EncodePath(mod)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		doc, err := proxy.GetDoc(r.Context(), mod, ver)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		err = tt.Lookup("index.html").Execute(w, map[string]interface{}{
			"index": false,
			"data":  doc,
		})
	}
}
