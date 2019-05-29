// Package fetch provides utilities to making authenticated
// fetches to the GOPROXY URL
package fetch

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
)

// Fetch makes a GET request to the given URL. It also appends
// an authentication token if GCP_SERVERLESS env is set to true.
// In the future, this should support any auth header and not just GCP.
func Fetch(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if os.Getenv("GCP_SERVERLESS") == "true" {
		u, err := neturl.Parse(url)
		if err != nil {
			return nil, fmt.Errorf("could not parse given url '%v': %v", url, err)
		}
		tok, err := getToken(u.Hostname())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	return http.DefaultClient.Do(req)
}

func getToken(audience string) (string, error) {
	url := "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Set("audience", audience)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata-Flavor", "Google")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not request creds: %v", err)
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read creds: %v", err)
	}

	return strings.TrimSpace(string(bts)), nil
}
