package asynqcodegen

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var (
	ErrFormat   = errors.New("asynq-codegen: could not format generated file")
	ErrGenerate = errors.New("asynq-codegen")
	ErrWrite    = errors.New("asynq-codegen: could not write to file")
)

// Generate writes a formatted Go file at the specified location.
func Generate(inspectResult *InspectResult, destination string) error {
	gen, err := Render(xrvTpl, inspectResult, nil)
	if err != nil {
		return errors.Join(ErrGenerate, err)
	}

	//nolint:gosec // It is fine to relax permissions of a generated file.
	if err := os.WriteFile(destination, gen, 0o644); err != nil {
		return errors.Join(ErrWrite, err)
	}

	if err := exec.CommandContext(context.Background(), "gofmt", "-w", destination).Run(); err != nil {
		return errors.Join(ErrFormat, err)
	}

	return nil
}

// Inspect parses Go code in the specified directory and returns an [InspectResult] that can be passed to [Generate].
func Inspect(cwd string) (*InspectResult, error) {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(
		fset,
		cwd,
		func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_generated.go")
		},
		parser.ParseComments,
	)
	if err != nil {
		//nolint:forbidigo // Keep this error contextualised.
		return nil, fmt.Errorf("parsing %s: %w", cwd, err)
	}

	out := &InspectResult{
		PackageName: "",
		Comments:    []AsynqComment{},
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				ac, err := inspectStructAnnotation(decl, pkg)

				switch {
				case err == nil:
					out.PackageName = pkg.Name
					out.Comments = append(out.Comments, *ac)
				case !errors.Is(err, errNoDirective):
					filename := fset.Position(file.Pos()).Filename
					line := fset.Position(file.Pos()).Line

					//nolint:forbidigo // Keep this error contextualised.
					return nil, fmt.Errorf("invalid asynq-codegen struct found in %s:%d: %w", filename, line, err)
				}
			}
		}
	}

	// This is to keep the output deterministic for the sake of repeatable code generations, cleaner diffs, etc.
	slices.SortFunc(out.Comments, cmpAsynqComments)

	return out, nil
}

func inspectStructAnnotation(decl ast.Decl, pkg *ast.Package) (*AsynqComment, error) {
	gen, ok := decl.(*ast.GenDecl)
	if !ok || gen.Tok != token.TYPE || len(gen.Specs) == 0 {
		return nil, errNoDirective
	}

	for _, spec := range gen.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		if _, ok := ts.Type.(*ast.StructType); !ok {
			continue
		}

		out, err := NewAsynqComment(gen.Doc, pkg.Name, ts.Name.Name)
		if err != nil {
			return nil, err
		}

		return out, nil
	}

	return nil, errNoDirective
}

func cmpAsynqComments(a, b AsynqComment) int {
	switch {
	case a.StructName < b.StructName:
		return -1
	case a.StructName > b.StructName:
		return 1
	default:
		return 0
	}
}
