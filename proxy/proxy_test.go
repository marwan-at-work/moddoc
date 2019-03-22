package proxy

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"testing"
// )

// /*
// // Hello Enums shi
// const (
// 	// One comment
// 	Hello = "there"

// 	Five = 3
// 	// Then is hahaha
// 	Then = "hahaha"
// )
// */

// func TestFull(t *testing.T) {
// 	s := &service{"http://localhost:3000"}
// 	doc, err := s.GetDoc(context.Background(), "github.com/pkg/errors", "v0.8.1")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Fatal(doc.PackageDoc)
// }

// func TestProxy(t *testing.T) {
// 	s := &service{}
// 	doc, err := s.getGoDoc(context.Background(), "hello", []*file{{
// 		Name: "hello.go",
// 		Content: []byte(`// Package main is pretty cool
// package main

// import "fmt"
// import "context"

// // OneNE is three
// var OneNE = Shiz()[3]

// // Shiz returns some shiz
// func Shiz(x fmt.Stringer, y string) (string, error) {
// 	return "", nil
// }

// // J shiz
// type J struct {
// 	OK string // hmmm
// 	Then int // commetime
// }

// // Alright methodz
// func (j *J) Alright() {

// }

// // NewJ inininit
// func NewJ() *J {
// 	return &J{}
// }

// func main() {

// }

// 		`),
// 	}})
// 	if err != nil {
// 		panic(err)
// 	}
// 	bts, _ := json.MarshalIndent(doc, "", "\t")
// 	fmt.Printf("%s\n", bts)
// 	t.Fatal("ok")
// }
