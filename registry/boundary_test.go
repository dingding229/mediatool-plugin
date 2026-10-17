package registry

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// A source guard, NOT a proof about arbitrary transitive initialization or malicious code.
func TestNoApplicationImportsOrPluginInitSideEffects(t *testing.T) {
	err := filepath.WalkDir("..", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, e := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if e != nil {
			return e
		}
		for _, imp := range f.Imports {
			value, e := strconv.Unquote(imp.Path.Value)
			if e != nil {
				return e
			}
			if strings.Contains(value, "media-tool/") || strings.Contains(value, "mediatool/internal") || strings.Contains(value, "wazero") || strings.Contains(value, "pgx") || strings.Contains(value, "go-redis") {
				t.Errorf("forbidden application/database dependency %s: %s", path, value)
			}
		}
		for _, decl := range f.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "init" {
				t.Errorf("package init forbidden: %s", path)
			}
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}
			for _, sp := range gen.Specs {
				vs := sp.(*ast.ValueSpec)
				for _, value := range vs.Values {
					ast.Inspect(value, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							// Compile-time interface assertion (*T)(nil) is a conversion, not execution.
							if paren, ok := call.Fun.(*ast.ParenExpr); ok {
								if _, ok = paren.X.(*ast.StarExpr); ok && len(call.Args) == 1 {
									if id, ok := call.Args[0].(*ast.Ident); ok && id.Name == "nil" {
										return false
									}
								}
							}
							t.Errorf("top-level function evaluation requires explicit review: %s", path)
						}
						return true
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
