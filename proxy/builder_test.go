package proxy

import (
	"fmt"
	"testing"
)

var getRelativeDirTestCases = []struct {
	file, mod, relativeTo string // input
	dir                   string // output
	valid                 bool   // output
}{
	{
		file:       "github.com/one/two@v0.3.0/main.go",
		relativeTo: "",
		dir:        ".",
		valid:      true,
	},
	{
		file:       "github.com/one/two@v0.3.0/three/main.go",
		relativeTo: "",
		dir:        "three",
		valid:      true,
	},
	{
		file:       "github.com/one/two@v0.3.0/three/main.go",
		relativeTo: "three",
		dir:        ".",
		valid:      true,
	},
	{
		file:       "github.com/one/two@v0.3.0/three/main.go",
		relativeTo: "threefour",
		dir:        "three",
		valid:      false,
	},
	{
		file:       "github.com/one/two@v0.3.0/three/main.go",
		relativeTo: "three/four",
		dir:        "three",
		valid:      false,
	},
}

func TestRelativeDir(t *testing.T) {
	for idx, tc := range getRelativeDirTestCases {
		t.Run(fmt.Sprintf("test_%v", idx), func(t *testing.T) {
			dir, valid := getRelativeDir(tc.file, tc.relativeTo)
			if valid != tc.valid {
				t.Fatalf(
					"expected the validity of file %v reltiveTo %v to be %v but got %v",
					tc.file, tc.relativeTo, tc.valid, valid,
				)
			}
			if dir != tc.dir {
				t.Fatalf("expected dir to be %v but got %v", tc.dir, dir)
			}
		})

	}
}

var getDirTestCases = []struct {
	input, output string
}{
	{"go.uber.org/zap@v1.9.1/README.md", "."},
	{"go.uber.org/zap@v1.9.1/benchmarks/scenario_bench_test.go", "benchmarks"},
	{"go.uber.org/zap@v1.9.1/benchmarks/scenario_bench_test.go", "benchmarks"},
	{"go.uber.org/zap@v1.9.1/zaptest/observer/logged_entry_test.go", "zaptest/observer"},
}

func TestGetDir(t *testing.T) {
	for idx, tc := range getDirTestCases {
		t.Run(fmt.Sprintf("test_%v", idx), func(t *testing.T) {
			given := getDir(tc.input)
			if tc.output != given {
				t.Fatalf("expected the directory for %v to be %v but got %v", tc.input, tc.output, given)
			}
		})
	}
}
