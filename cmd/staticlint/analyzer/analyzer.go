//go:build ignore

// Сервис проверки правил анализаторов.
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Start функция запуска методов по проверке правил
func Start() {
	var ErrCheckAnalyzer = &analysis.Analyzer{
		Name: "increment18",
		Doc:  "Проверка правил для инкремента 18",
		Run:  run,
	}

	mychecks := []*analysis.Analyzer{
		ErrCheckAnalyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

	// добавляем анализаторы из staticcheck
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// добавляем анализаторы из simple
	for _, v := range simple.Analyzers {
		switch v.Analyzer.Name {
		case "S1001":
			mychecks = append(mychecks, v.Analyzer)
		case "S1003":
			mychecks = append(mychecks, v.Analyzer)
		case "S1008":
			mychecks = append(mychecks, v.Analyzer)
		default:
			continue
		}
	}

	// добавляем анализаторы из stylecheck
	for _, v := range stylecheck.Analyzers {
		switch v.Analyzer.Name {
		case "ST1001":
			mychecks = append(mychecks, v.Analyzer)
		case "ST1005":
			mychecks = append(mychecks, v.Analyzer)
		case "ST1013":
			mychecks = append(mychecks, v.Analyzer)
		default:
			continue
		}
	}

	// добавляем анализаторы из quickfix
	for _, v := range quickfix.Analyzers {
		switch v.Analyzer.Name {
		case "QF1001":
			mychecks = append(mychecks, v.Analyzer)
		case "QF1004":
			mychecks = append(mychecks, v.Analyzer)
		case "QF1010":
			mychecks = append(mychecks, v.Analyzer)
		default:
			continue
		}
	}

	multichecker.Main(
		mychecks...,
	)
}

// run функция, реализующая пользовательские проверки
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			nameMain := funcDecl.Name.Name
			if nameMain != "main" {
				return true
			}

			for _, stmt := range funcDecl.Body.List {
				exprStmt, ok := stmt.(*ast.ExprStmt)
				if !ok {
					return true
				}

				call, ok := exprStmt.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					funcName := fun.Sel.Name

					if funcName == "Exit" {
						pass.Reportf(fun.Pos(), "not allowed to call the Exit")
					}
				}
			}

			return true
		})
	}

	return nil, nil
}
