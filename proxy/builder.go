package proxy

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/doc"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	proxydoc "marwan.io/moddoc/doc"
	"marwan.io/moddoc/gocopy/modfile"
	"marwan.io/moddoc/gocopy/module"
)

type builder struct {
	fset     *token.FileSet
	examples []*doc.Example
	mods     []*modFile
}

func (b *builder) getGoDoc(ctx context.Context, mod, ver, subpkg string, files []*file) (*proxydoc.Documentation, error) {
	b.fset = token.NewFileSet()
	mp := map[string]*ast.File{}
	dirMap := map[string]struct{}{}
	pkgName := ""
	// TODO: parse sub directories to get synopsis
	pkgFiles := []*proxydoc.File{}
	testFiles := []*ast.File{}
	for _, f := range files {
		if filepath.Base(f.Name) == "go.mod" {
			modf, err := modfile.Parse("go.mod", f.Content, nil)
			if err != nil {
				return nil, err
			}
			b.mods = append(b.mods, &modFile{path: filepath.Dir(f.Name), file: modf})
			continue
		}
		if filepath.Ext(f.Name) != ".go" {
			continue
		}
		dir, valid := getRelativeDir(f.Name, subpkg)
		if !valid {
			continue
		}
		if dir != "." {
			dirMap[dir] = struct{}{}
			continue
		}
		astFile, err := parser.ParseFile(b.fset, f.Name, f.Content, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(f.Name, "_test.go") {
			testFiles = append(testFiles, astFile)
			continue
		}
		mp[f.Name] = astFile
		if pkgName == "" {
			pkgName = astFile.Name.String()
		}
		pkgFiles = append(pkgFiles, &proxydoc.File{Name: filepath.Base(f.Name)})
	}
	b.examples = doc.Examples(testFiles...)
	astPkg := &ast.Package{Name: mod, Files: mp}
	dpkg := doc.New(astPkg, mod, doc.Mode(0))
	var d proxydoc.Documentation
	d.PackageName = pkgName
	var sb strings.Builder
	doc.ToHTML(&sb, dpkg.Doc, nil)
	d.PackageDoc = template.HTML(sb.String())
	d.ImportPath, _ = module.DecodePath(mod)
	d.Constants = b.getConsts(dpkg.Consts)
	d.Variables = b.getConsts(dpkg.Vars)
	d.Funcs = b.getFuncs(dpkg.Funcs, "")
	d.Types = b.getTypes(dpkg.Types)
	d.Files = pkgFiles
	d.Examples = b.getExamples("")
	d.ModuleVersion = ver
	for subDir := range dirMap {
		d.Subdirs = append(d.Subdirs, &proxydoc.Subdir{
			Name:     subDir,
			Synopsis: getSynopsis(subDir, files),
			Link:     filepath.Join("/", d.ImportPath, subDir, "@v", d.ModuleVersion),
		})
	}
	sort.Slice(d.Subdirs, func(i, j int) bool {
		return d.Subdirs[i].Name < d.Subdirs[j].Name
	})

	if len(b.mods) > 0 {
		modf := b.getClosestModFile(mod)
		d.GoMod = b.getMod(modf.file)
	}

	d.NavLinks = []string{"Index"}
	if len(d.Examples) > 0 {
		d.NavLinks = append(d.NavLinks, "Examples")
	}
	if len(d.Files) > 0 {
		d.NavLinks = append(d.NavLinks, "Files")
	}
	if len(d.GoMod) > 0 {
		d.NavLinks = append(d.NavLinks, "Go.mod")
	}
	if len(d.Subdirs) > 0 {
		d.NavLinks = append(d.NavLinks, "Directories")
	}

	return &d, nil
}

func (b *builder) getMod(mod *modfile.File) template.HTML {
	mp := map[string]string{}
	for _, req := range mod.Require {
		mp[req.Mod.Path] = fmt.Sprintf(`/%s/@v/%s`, req.Mod.Path, req.Mod.Version)
	}
	for _, rep := range mod.Replace {
		mp[rep.New.Path] = fmt.Sprintf(`/%s/@v/%s`, rep.New.Path, rep.New.Version)
	}
	return template.HTML(modfile.FormatHTML(mod.Syntax, mp))
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
	t.Funcs = b.getFuncs(typ.Funcs, "")
	t.Methods = b.getFuncs(typ.Methods, t.Name)
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
	t.Doc = template.HTML(docStr.String())
	t.Constants = b.getConsts(typ.Consts)
	t.Variables = b.getConsts(typ.Vars)
	t.Examples = b.getExamples(t.Name)

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

func (b *builder) getFuncs(funcs []*doc.Func, typeName string) []*proxydoc.Func {
	res := []*proxydoc.Func{}
	for _, f := range funcs {
		ff := b.getFunc(f, typeName)
		res = append(res, ff)
	}
	return res
}

func (b *builder) getFunc(f *doc.Func, typeName string) *proxydoc.Func {
	var df proxydoc.Func
	df.ID = f.Name
	if typeName != "" {
		df.ID = typeName + "." + f.Name
	}
	df.Name = f.Name
	var docBuilder strings.Builder
	doc.ToHTML(&docBuilder, f.Doc, nil)
	df.Doc = template.HTML(docBuilder.String())
	// df.Signature = &proxydoc.FunctionSignature{} //TODO: make receiver/args/returns clickable.
	var sb strings.Builder
	err := format.Node(&sb, b.fset, f.Decl)
	if err != nil {
		fmt.Println("could not format function signature", err)
	}
	df.SignatureString = sb.String()
	df.MethodReceiverString = f.Recv
	examplePrefix := df.Name
	if typeName != "" {
		examplePrefix += "_" + typeName
	}
	df.Examples = b.getExamples(examplePrefix)

	return &df
}

func (b *builder) getConsts(cc []*doc.Value) []*proxydoc.Value {
	vals := []*proxydoc.Value{}
	for _, c := range cc {
		var docBuilder strings.Builder
		doc.ToHTML(&docBuilder, c.Doc, nil)
		val := &proxydoc.Value{
			IsGroup: len(c.Names) > 1,
			Doc:     template.HTML(docBuilder.String()),
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
				newV.Doc = template.HTML(spec.Doc.Text())
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
func getDir(zipPath string) string {
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

func getRelativeDir(file, relativeTo string) (dir string, valid bool) {
	dir = getDir(file)
	// if module root, return all dirs so that we can capture sub directories
	if relativeTo == "" {
		return dir, true
	}
	// if outside of the relative path, it should not exist
	if !strings.HasPrefix(dir+"/", relativeTo+"/") {
		return dir, false
	}
	if dir == relativeTo {
		return ".", true
	}
	// todo: ensure robust.
	return dir[len(relativeTo)+1:], true
}

func (b *builder) getExamples(name string) []*proxydoc.Example {
	var docs []*proxydoc.Example
	for _, e := range b.examples {
		if !strings.HasPrefix(e.Name, name) {
			continue
		}
		n := e.Name[len(name):]
		if n != "" {
			if i := strings.LastIndex(n, "_"); i != 0 {
				continue
			}
			n = n[1:]
			if startsWithUppercase(n) {
				continue
			}
			n = strings.Title(n)
		}

		var codeBuilder strings.Builder
		var nn interface{}
		if _, ok := e.Code.(*ast.File); ok {
			nn = e.Play
		} else {
			nn = &printer.CommentedNode{Node: e.Code, Comments: e.Comments}
		}
		err := (&printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}).Fprint(&codeBuilder, b.fset, nn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		code, output := fmtExampleCode(codeBuilder.String(), e.Output)

		docs = append(docs, &proxydoc.Example{
			ID:     "Example" + name + "--" + n,
			Name:   n,
			Doc:    e.Doc,
			Code:   code,
			Output: output,
			// Play:   play,
		})
	}
	return docs
}

type modFile struct {
	path string
	file *modfile.File
}

func (b *builder) getClosestModFile(dir string) *modFile {
	if len(b.mods) == 1 {
		return b.mods[0]
	}

	var longest *modFile
	for _, m := range b.mods {
		if longest != nil && len(longest.path) > len(m.path) {
			continue
		}
		if strings.HasPrefix(dir, m.path) {
			longest = m
		}
	}
	return longest
}

var exampleOutputRx = regexp.MustCompile(`(?i)//[[:space:]]*output:`)

func fmtExampleCode(s, output string) (string, string) {
	buf := []byte(s)

	// additional formatting if this is a function body
	if i := len(buf); i >= 2 && buf[0] == '{' && buf[i-1] == '}' {
		// remove surrounding braces
		buf = buf[1 : i-1]
		// unindent
		buf = bytes.Replace(buf, []byte("\n    "), []byte("\n"), -1)
		// remove output comment
		if j := exampleOutputRx.FindIndex(buf); j != nil {
			buf = bytes.TrimSpace(buf[:j[0]])
		}
	} else {
		// drop output, as the output comment will appear in the code
		output = ""
	}
	return string(buf), output
}

func startsWithUppercase(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
