package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"marwan.io/moddoc/fetch"
)

type indexResponse struct {
	Path      string
	Version   string
	Timestamp time.Time
}

func index(ctx context.Context) ([]*moduleIndex, error) {
	url := "https://index.golang.org/index"
	resp, err := fetch.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %v", resp.Status)
	}

	mp := map[string][]string{}
	dec := json.NewDecoder(resp.Body)
	for dec.More() {
		var ir indexResponse
		err := dec.Decode(&ir)
		if err != nil {
			return nil, err
		}
		mp[ir.Path] = append(mp[ir.Path], ir.Version)
	}

	mods := []*moduleIndex{}

	for mod, vers := range mp {
		mods = append(mods, &moduleIndex{
			Module:   mod,
			Versions: vers,
			Latest:   latestVer(vers),
		})
	}

	return mods, nil
}
