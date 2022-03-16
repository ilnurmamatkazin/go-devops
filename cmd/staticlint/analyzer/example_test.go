package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func ExampleStart() {
	// Для обработки пользовательских правил необходимо объявить переменную типа analysis.Analyzer
	// В объявленной переменной необходимо указать в ней функцию обработки пользовательких правил run.
	var ErrCheckAnalyzer = &analysis.Analyzer{
		Name: "increment18",
		Doc:  "Проверка правил для инкремента 18",
		Run:  _run,
	}

	// Формируем слайс из всех необходимых нам анализаторов.
	// Первым указываем самописный анализатор с пользовательскими правилами.
	// Далее идет список стандартных анализаторов.
	mychecks := []*analysis.Analyzer{
		ErrCheckAnalyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		//... продолжаем дополнять список необходимыми анализаторами
	}

	// Добавляем все анализаторы из пакета staticcheck
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	// Добавляем часть анализаторов из пакета simple
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

	// Добавляем часть анализаторов из пакета stylecheck
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

	// Добавляем часть анализаторов из пакета quickfix
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

	// Помещаем итоговый слайс с анализаторами в multichecker
	multichecker.Main(
		mychecks...,
	)
}

// Функция, реализующая логику проверки кода на соответствие пользовательским требованиям.
func _run(pass *analysis.Pass) (interface{}, error) {
	// проходим по всем файлам проекта
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			// определяем, что узел является декларацией функции
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// проверяем, что декларация функции соответсвует функии main
			nameMain := funcDecl.Name.Name
			if nameMain != "main" {
				return true
			}

			// проходим по списку узлов функции main
			for _, stmt := range funcDecl.Body.List {
				// определяем, что узел является автономным выражением в списке операторов
				exprStmt, ok := stmt.(*ast.ExprStmt)
				if !ok {
					return true
				}

				// определяем, что узел представляет собой выражение, за которым следует список аргументов
				call, ok := exprStmt.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				// определяем, что узел является выражением, за которым следует селектор
				// и определяем его имя
				if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
					funcName := fun.Sel.Name

					// проверяем, что имя селектора соответсвует оператору Exit и сообщаем об этом
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
