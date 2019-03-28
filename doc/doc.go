package doc

import (
	"html/template"
)

// Documentation is the data structure
// that represents a full module page
type Documentation struct {
	PackageName   string
	ModuleVersion string
	Versions      []string
	ModuleRoot    string
	ImportPath    string
	PackageDoc    template.HTML
	Constants     []*Value
	Variables     []*Value
	Funcs         []*Func
	Types         []*Type
	Files         []*File
	Subdirs       []*Subdir
	NavLinks      []string
}

// Value represents one or a group of constants/variables
type Value struct {
	SignatureString string
	Name            string
	Value           string
	Type            string
	Doc             template.HTML
	IsGroup         bool
	Values          []*Value
}

// Func represents a function or a method
type Func struct {
	ID   string // Name for funcs; TypeName+FuncName for type methods.
	Name string
	// Signature       *FunctionSignature //TODO: later
	SignatureString string
	Doc             template.HTML
	// MethodReceiver  *MethodReceiver // TODO: later
	MethodReceiverString string
	Examples             []*Example
}

// FunctionSignature represents a function or method signature
type FunctionSignature struct {
	Arguments []*Argument
	Returns   []*Argument
}

// Argument is either an input or return name/type
type Argument struct {
	Name       string
	Type       string
	IsVariadic bool
}

// Example represents a type or function example
type Example struct {
	Name   string
	Code   string
	Output string
}

// Type represents a type declaration
type Type struct {
	Name            string
	Doc             template.HTML
	Type            string
	SignatureString string
	Fields          []*Field
	Examples        []*Example
	Methods         []*Func
	Funcs           []*Func
	Constants       []*Value
	Variables       []*Value
}

// Field is a struct filed
type Field struct {
	Name      string
	Type      string
	Doc       string
	StructTag string
}

// MethodReceiver is a method receiver
// that belongs to a type. The struct
// assumes the parent knows what the type name is.
type MethodReceiver struct {
	Name      string
	IsPointer bool
}

// File represents a go file inside a package
type File struct {
	Name string
	// Future: link
}

// Subdir represents a potential sub package.
// Caller assumes they know how to link to a module
// to one of its sub-directories.
type Subdir struct {
	Name     string
	Synopsis string
	Link     string
}
