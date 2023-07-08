package gofeatureenvy

import (
	"fmt"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

const doc = "gofeatureenvy "

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "gofeatureenvy",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
}

const (
	few = 5
)

func run(pass *analysis.Pass) (any, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)

	type StructMetrics struct {
		pos              token.Pos
		localAccess      uint
		foreignAccess    uint
		foreignStructMap map[types.Type]struct{}
	}
	structMap := make(map[*types.TypeName]*StructMetrics)
	for _, f := range s.SrcFuncs {
		if f.Signature == nil || f.Signature.Recv() == nil {
			continue
		}

		derefedType := f.Signature.Recv().Type()
		for pointerType, ok := derefedType.(*types.Pointer); ok; pointerType, ok = derefedType.(*types.Pointer) {
			derefedType = pointerType.Elem()
		}
		namedType, ok := derefedType.(*types.Named)
		if !ok || namedType == nil {
			continue
		}

		structType, ok := namedType.Underlying().(*types.Struct)
		if !ok || structType == nil {
			continue
		}

		structMetrics, ok := structMap[namedType.Obj()]
		if !ok {
			structMetrics = &StructMetrics{
				pos:              namedType.Obj().Pos(),
				localAccess:      0,
				foreignAccess:    0,
				foreignStructMap: map[types.Type]struct{}{},
			}
			structMap[namedType.Obj()] = structMetrics
		}

		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				switch instr := instr.(type) {
				case *ssa.FieldAddr:
					if instr.X == nil {
						continue
					}

					if instr.X.Type() == namedType {
						structMetrics.localAccess++
					} else {
						structMetrics.foreignAccess++
						structMetrics.foreignStructMap[instr.X.Type()] = struct{}{}
					}
				case *ssa.Field:
					if instr.X == nil {
						continue
					}

					if instr.X.Type() == namedType {
						structMetrics.localAccess++
					} else {
						structMetrics.foreignAccess++
						structMetrics.foreignStructMap[instr.X.Type()] = struct{}{}
					}
				}
			}
		}
	}

	for _, structMetrics := range structMap {
		atfd := structMetrics.foreignAccess
		laa := float64(structMetrics.localAccess) / float64(structMetrics.foreignAccess+structMetrics.localAccess)
		fdp := len(structMetrics.foreignStructMap)

		fmt.Printf("atfd: %d, laa: %f, fdp: %d\n", atfd, laa, fdp)
		if atfd > few && laa < 1.0/3 && fdp <= few {
			pass.Reportf(structMetrics.pos, "feature envy")
		}
	}
	return nil, nil
}
