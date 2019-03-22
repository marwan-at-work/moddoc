package proxy

import (
	"context"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"

	proxydoc "marwan.io/moddoc/doc"
	"marwan.io/moddoc/gocopy/module"
)

type builder struct {
	fset *token.FileSet
}

func (b *builder) getGoDoc(ctx context.Context, mod, subpkg string, files []*file) (*proxydoc.Documentation, error) {
	b.fset = token.NewFileSet()
	mp := map[string]*ast.File{}
	dirMap := map[string]struct{}{}
	pkgName := ""
	// TODO: parse sub directories to get synopsis
	pkgFiles := []*proxydoc.File{}
	for _, f := range files {
		dir, valid := getRelativeDir(f.Name, mod, subpkg)
		if !valid {
			continue
		}
		if dir != "." {
			dirMap[dir] = struct{}{}
			continue
		}
		if filepath.Ext(f.Name) != ".go" || strings.HasSuffix(f.Name, "_test.go") {
			continue
		}
		astFile, err := parser.ParseFile(b.fset, f.Name, f.Content, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		mp[f.Name] = astFile
		if pkgName == "" {
			pkgName = astFile.Name.String()
		}
		pkgFiles = append(pkgFiles, &proxydoc.File{Name: filepath.Base(f.Name)})
	}
	astPkg := &ast.Package{Name: mod, Files: mp}
	dpkg := doc.New(astPkg, mod, doc.Mode(0))
	var d proxydoc.Documentation
	d.PackageName = pkgName
	var sb strings.Builder
	doc.ToHTML(&sb, dpkg.Doc, nil)
	d.PackageDoc = sb.String()
	d.ImportPath, _ = module.DecodePath(mod)
	d.Constants = b.getConsts(dpkg.Consts)
	d.Variables = b.getConsts(dpkg.Vars)
	d.Funcs = b.getFuncs(dpkg.Funcs)
	d.Types = b.getTypes(dpkg.Types)
	d.Files = pkgFiles

	for subDir := range dirMap {
		d.Subdirs = append(d.Subdirs, &proxydoc.Subdir{
			Name:     subDir,
			Synopsis: getSynopsis(subDir, files),
		})
	}
	sort.Slice(d.Subdirs, func(i, j int) bool {
		return d.Subdirs[i].Name < d.Subdirs[j].Name
	})
	return &d, nil
}

func getSynopsis(subDir string, files []*file) string {
	for _, f := range files {
		dir := filepath.Dir(f.Name)
		fileInDir := strings.HasSuffix(dir, subDir)
		ext := filepath.Ext(f.Name)
		if fileInDir && ext == ".go" {
			astFile, err := parser.ParseFile(token.NewFileSet(), f.Name, f.Content, parser.ParseComments)
			if err != nil {
				fmt.Println(err)
				// TODO: handle
			}
			if astFile.Doc != nil {
				return synopsis(astFile.Doc.Text())
			}
		}
	}

	return ""
}

func (b *builder) getTypes(types []*doc.Type) []*proxydoc.Type {
	tt := []*proxydoc.Type{}
	for _, t := range types {
		tt = append(tt, b.getType(t))
	}
	return tt
}

func (b *builder) getType(typ *doc.Type) *proxydoc.Type {
	var t proxydoc.Type
	t.Name = typ.Name
	t.Funcs = b.getFuncs(typ.Funcs)
	t.Methods = b.getFuncs(typ.Methods)
	spec := typ.Decl.Specs[0].(*ast.TypeSpec)
	// todo: must use original FileSet for struct inline comments
	if structType, ok := spec.Type.(*ast.StructType); ok {
		t.Type = "struct"
		t.Fields = b.getFields(structType)
	} else {
		var sb strings.Builder
		format.Node(&sb, b.fset, spec)
		t.Type = sb.String()
	}
	var sb strings.Builder
	format.Node(&sb, b.fset, typ.Decl)
	t.SignatureString = sb.String()
	var docStr strings.Builder
	doc.ToHTML(&docStr, typ.Doc, nil)
	t.Doc = docStr.String()
	t.Constants = b.getConsts(typ.Consts)
	t.Variables = b.getConsts(typ.Vars)

	return &t
}

func (b *builder) getFields(st *ast.StructType) []*proxydoc.Field {
	fields := []*proxydoc.Field{}
	for _, f := range st.Fields.List {
		field := b.getField(f)
		if field != nil {
			fields = append(fields, field)
		}
	}
	return fields
}

// TODO: support embedded structs
// TODO: support inline structs
// TODO: support multi named fields
func (b *builder) getField(f *ast.Field) *proxydoc.Field {
	var df proxydoc.Field
	if len(f.Names) == 0 {
		return nil
	}
	df.Name = f.Names[0].Name
	var sb strings.Builder
	format.Node(&sb, b.fset, f.Type)
	df.Type = sb.String()
	df.Doc = f.Doc.Text()
	if f.Tag != nil {
		df.StructTag = f.Tag.Value
	}

	return &df
}

func (b *builder) getFuncs(funcs []*doc.Func) []*proxydoc.Func {
	res := []*proxydoc.Func{}
	for _, f := range funcs {
		ff := b.getFunc(f)
		res = append(res, ff)
	}
	return res
}

func (b *builder) getFunc(f *doc.Func) *proxydoc.Func {
	var df proxydoc.Func
	df.Name = f.Name
	var docBuilder strings.Builder
	doc.ToHTML(&docBuilder, f.Doc, nil)
	df.Doc = docBuilder.String()
	// df.Signature = &proxydoc.FunctionSignature{} //TODO: make receiver/args/returns clickable.
	var sb strings.Builder
	err := format.Node(&sb, b.fset, f.Decl)
	if err != nil {
		fmt.Println("could not format function signature", err)
	}
	df.SignatureString = sb.String()
	df.MethodReceiverString = f.Recv

	return &df
}

func (b *builder) getConsts(cc []*doc.Value) []*proxydoc.Value {
	vals := []*proxydoc.Value{}
	for _, c := range cc {
		var docBuilder strings.Builder
		doc.ToHTML(&docBuilder, c.Doc, nil)
		val := &proxydoc.Value{
			IsGroup: len(c.Names) > 1,
			Doc:     docBuilder.String(),
		}
		if val.IsGroup {
			for idx, n := range c.Names {
				newV := &proxydoc.Value{
					IsGroup: false,
					Name:    n,
				}
				spec, ok := c.Decl.Specs[idx].(*ast.ValueSpec)
				if !ok {
					fmt.Printf("unrecognized group spec type: %T\n", c.Decl.Specs[idx])
					return vals
				}
				newV.Doc = spec.Doc.Text()
				b.populateConstantsValueAndType(newV, spec)
				val.Values = append(val.Values, newV)
			}
		} else {
			val.Name = c.Names[0]
			spec, ok := c.Decl.Specs[0].(*ast.ValueSpec)
			if !ok {
				fmt.Printf("unrecognized spec type: %T\n", c.Decl.Specs[0])
				return vals
			}
			b.populateConstantsValueAndType(val, spec)
		}
		var sb strings.Builder
		format.Node(&sb, b.fset, c.Decl)
		val.SignatureString = sb.String()
		vals = append(vals, val)
	}
	return vals
}

func (b *builder) populateConstantsValueAndType(v *proxydoc.Value, spec *ast.ValueSpec) {
	if len(spec.Values) == 1 {
		v.Value = b.getValueFromSpec(v, spec)
	}
	v.Type = b.getTypeFromSpec(spec)
}

func (b *builder) getValueFromSpec(v *proxydoc.Value, spec *ast.ValueSpec) string {
	var sb strings.Builder
	err := format.Node(&sb, b.fset, spec.Values[0])
	if err != nil {
		fmt.Println("value formatting error", err)
	}
	return sb.String()
}

func (b *builder) getTypeFromSpec(spec *ast.ValueSpec) string {
	if spec.Type != nil {
		var sb strings.Builder
		format.Node(&sb, b.fset, spec.Type)
		return sb.String()
	}
	return ""
}

// TODO: might be better (or not) to regex against (.+@[^/]+)/(.+)
func getDir(zipPath, mod string) string {
	idx := strings.Index(zipPath, "@")
	if idx == -1 {
		idx = 0
	}
	zipPath = zipPath[idx:]
	idx = strings.Index(zipPath, "/")
	if idx == -1 {
		idx = 0
	}
	return filepath.Dir(zipPath[idx+1:])
}

func getRelativeDir(file, mod, relativeTo string) (dir string, valid bool) {
	dir = getDir(file, mod)
	// if module root, return all dirs so thatwe can capture sub directories
	if relativeTo == "" {
		return dir, true
	}
	// if outside of the relative path, it should not exist
	idx := strings.Index(dir, relativeTo)
	if idx == -1 {
		return dir, false
	}
	if dir == relativeTo {
		return ".", true
	}
	// todo: ensure robust.
	return dir[idx+len(relativeTo)+1:], true
}
