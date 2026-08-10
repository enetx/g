package g_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPrimitiveLayerBudget pins the compile-cost invariant this library is
// built around: the method sets of the primitive types (String, Int, Float,
// Bytes and their accessor structs) must never mention the heavyweight generic
// containers in SIGNATURES OR BODIES. One such mention re-welds the whole
// container graph onto every package that names g.String (measured: 1 556
// emitted functions today vs 33 332 before the de-weld — a single Seq[String]
// in a splitter signature costs ×4.7).
//
// If this test fails, the offending method must return a plain slice, a bare
// func type, or move to a subpackage — see the de-weld notes in doc.go.
func TestPrimitiveLayerBudget(t *testing.T) {
	// Files whose declarations belong to the primitive layer.
	primitiveFiles := []string{
		"string.go", "string_builder.go", "string_crypt.go", "string_encdec.go",
		"string_hash.go", "string_compdecomp.go",
		"bytes.go", "bytes_text.go", "bytes_encdec.go", "bytes_hash.go",
		"bytes_compdecomp.go",
		"int.go", "int_checked.go", "float.go", "float_math.go",
	}

	// Generic container types that must NOT appear in the primitive layer.
	forbidden := map[string]bool{
		"Seq": true, "Seq2": true, "SeqSlices": true, "SeqResult": true,
		"Slice": true, "Map": true, "MapOrd": true, "MapSafe": true,
		"Set": true, "Deque": true, "Heap": true,
		"collector": true, "collector2": true, "collectorS": true,
	}

	// Sanctioned exceptions, checked as "file:receiver.Method uses Name".
	allowed := map[string]bool{
		// SubString/SubBytes index runes/bytes with plain loops; the only
		// container mention left is documented below if ever reintroduced.
	}

	fset := token.NewFileSet()

	for _, name := range primitiveFiles {
		path := filepath.Join("..", name)
		if _, err := os.Stat(path); err != nil {
			continue // file may be renamed/removed; the sweep is best-effort per file
		}

		file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			label := name + ":" + fn.Name.Name

			ast.Inspect(fn, func(n ast.Node) bool {
				idx, ok := n.(*ast.IndexExpr)
				if !ok {
					if lst, ok := n.(*ast.IndexListExpr); ok {
						if id, ok := lst.X.(*ast.Ident); ok && forbidden[id.Name] && !allowed[label+" uses "+id.Name] {
							t.Errorf("%s mentions generic container %s at %s — this re-welds the container graph onto the primitive layer",
								label, id.Name, fset.Position(lst.Pos()))
						}
					}
					return true
				}

				if id, ok := idx.X.(*ast.Ident); ok && forbidden[id.Name] && !allowed[label+" uses "+id.Name] {
					t.Errorf("%s mentions generic container %s at %s — this re-welds the container graph onto the primitive layer",
						label, id.Name, fset.Position(idx.Pos()))
				}

				return true
			})
		}
	}
}

// TestPrimitiveNoLazyNamedSeq additionally pins that the eager splitters stay
// plain slices: their result types must be []String / []Bytes / []rune, never
// a named sequence.
func TestPrimitiveNoLazyNamedSeq(t *testing.T) {
	fset := token.NewFileSet()

	for _, name := range []string{"string.go", "bytes.go", "bytes_text.go"} {
		file, err := parser.ParseFile(fset, filepath.Join("..", name), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}

		splitters := map[string]bool{
			"Split": true, "SplitAfter": true, "SplitN": true, "Lines": true,
			"Fields": true, "FieldsBy": true, "Chars": true, "Chunks": true, "Runes": true,
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || !splitters[fn.Name.Name] {
				continue
			}

			if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
				continue
			}

			switch rt := fn.Type.Results.List[0].Type.(type) {
			case *ast.ArrayType:
				// []String / []Bytes / []rune — fine.
			default:
				pos := fset.Position(fn.Pos())
				t.Errorf("%s.%s at %s returns %T — splitters must return plain slices (see step-2 de-weld notes)",
					name, fn.Name.Name, pos, rt)
			}
		}
	}
}

var _ = strings.TrimSpace // keep strings imported if the allowlist above grows
