package proxy

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	proxydoc "marwan.io/moddoc/doc"
	"marwan.io/moddoc/gocopy/module"
)

// Service can return a valid godoc
type Service interface {
	GetDoc(ctx context.Context, mod, ver string) (*proxydoc.Documentation, error)
}

// NewService returns a valid service based on a GOPROXY
func NewService(url string) Service {
	return &service{url: strings.TrimSuffix(url, "/")}
}

type service struct {
	url string
}

// GetProxyDir from GOPROXY
func (s *service) GetDoc(ctx context.Context, mod, ver string) (*proxydoc.Documentation, error) {
	dir, fileName, subpkg, err := s.makeZip(ctx, mod, ver)
	if err != nil {
		return nil, fmt.Errorf("could not make zip: %v", err)
	}
	modRoot := mod
	if subpkg != "" {
		rootIdx := len(mod) - len(subpkg)
		modRoot = strings.TrimSuffix(mod[0:rootIdx], "/")
	}
	versCh := s.getVersions(ctx, modRoot)
	defer os.RemoveAll(dir)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	zipReader, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, err
	}
	files := []*file{}
	// TODO: parse sub directories to get synopsis
	for _, f := range zipReader.File {
		var fl file
		fl.Name = f.Name
		rdr, err := f.Open()
		if err != nil {
			return nil, err
		}
		bts, err := ioutil.ReadAll(rdr)
		if err != nil {
			return nil, err
		}
		fl.Content = bts
		files = append(files, &fl)
	}

	bldr := &builder{}
	proxyDoc, err := bldr.getGoDoc(ctx, mod, ver, subpkg, files)
	proxyDoc.ModuleRoot, _ = module.DecodePath(modRoot)
	proxyDoc.Versions = <-versCh
	return proxyDoc, err
}

type file struct {
	Name    string
	Content []byte
}

func (s *service) getVersions(ctx context.Context, mod string) chan []string {
	ch := make(chan []string, 1)
	go func() {
		defer close(ch)
		resp, err := s.fetch(ctx, mod, "list", "")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return
		}
		bts, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		ch <- strings.Split(string(bts), "\n")
	}()
	return ch
}

func (s *service) makeZip(ctx context.Context, mod, ver string) (string, string, string, error) {
	dir, err := ioutil.TempDir("", strings.Replace(mod, "/", "_", -1)+ver)
	if err != nil {
		return dir, "", "", err
	}
	path := mod
	var resp *http.Response
	var subdir string
	for {
		if path == "." {
			return "", "", "", fmt.Errorf("invalid path: %v", mod)
		}
		resp, err = s.fetch(ctx, path, ver, ".zip")
		if err != nil {
			return "", "", "", err
		}
		if resp.StatusCode == 200 {
			break
		}
		path = filepath.Dir(path)
		subdir = mod[len(path)+1:]
	}
	file := filepath.Join(dir, "source.zip")
	f, err := os.Create(file)
	if err != nil {
		return dir, file, "", err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return dir, file, subdir, err
}

func (s *service) fetch(ctx context.Context, mod, ver, ext string) (*http.Response, error) {
	req, err := http.NewRequest("GET", s.url+"/"+mod+"/@v/"+ver+ext, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return http.DefaultClient.Do(req)
}
