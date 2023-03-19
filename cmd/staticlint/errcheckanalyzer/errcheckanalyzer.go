package errcheckanalyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Pass struct {
	// отобразим здесь только важные поля
	Fset         *token.FileSet // информация о позиции токенов
	Files        []*ast.File    // AST для каждого файла
	OtherFiles   []string       // имена файлов не на Go в пакете
	IgnoredFiles []string       // имена игнорируемых исходных файлов в пакете
	Pkg          *types.Package // информация о типах пакета
	TypesInfo    *types.Info    // информация о типах в AST
}

// ErrCheckAnalyzer init analyzer to check os.Exit() string in main func of main pkg
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck18",
	Doc:  "check for os.Exit() string in main func of main pkg",
	Run:  run,
}

// FuncMainEnd store num of last line in main func
var FuncMainEnd int

/*
	func run(pass *analysis.Pass) (interface{}, error) {
	    // реализация будет ниже
	    return nil, nil
	}
*/

// errorType is error type
var errorType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

// isErrorType checks error type
func isErrorType(t types.Type) bool {
	return types.Implements(t, errorType)
}

// resultErrors возвращает булев массив со значениями true,
// если тип i-го возвращаемого значения соответствует ошибке.
func resultErrors(pass *analysis.Pass, call *ast.CallExpr) []bool {
	switch t := pass.TypesInfo.Types[call].Type.(type) {
	case *types.Named: // возвращается значение
		return []bool{isErrorType(t)}
	case *types.Pointer: // возвращается указатель
		return []bool{isErrorType(t)}
	case *types.Tuple: // возвращается несколько значений
		s := make([]bool, t.Len())
		for i := 0; i < t.Len(); i++ {
			switch mt := t.At(i).Type().(type) {
			case *types.Named:
				s[i] = isErrorType(mt)
			case *types.Pointer:
				s[i] = isErrorType(mt)
			}
		}
		return s
	}
	return []bool{false}
}

// isReturnError возвращает true, если среди возвращаемых значений есть ошибка.
func isReturnError(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, isError := range resultErrors(pass, call) {
		if isError {
			return true
		}
	}
	return false
}

// run analysis
func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		// проверяем, что выражение представляет собой вызов функции,
		// у которой возвращаемая ошибка никак не обрабатывается
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isReturnError(pass, call) {
				pass.Reportf(x.Pos(), "expression returns unchecked error")
			}
		}
	}

	FuncMainsearchEnd := func(x *ast.FuncDecl) /*(endPos token.Pos)*/ {
		//pass.Reportf(x.Pos(), x.Name.String())
		//pass.Reportf(x.End(), "end of func")
		//fmt.Println(pass.Fset.Position(x.End()).Line)
		FuncMainEnd = (pass.Fset.Position(x.End()).Line)
	}

	funcIdent := func(x *ast.Ident) {
		if pass.Fset.Position(x.Pos()).Line <= FuncMainEnd {
			pass.Reportf(x.Pos(), "os."+x.Name+" is forbidden to call from main function of main pkg ")
		}
	}

	//n.(*ast.FuncType)
	lastIdent := ""
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			// функцией ast.Inspect проходим по всем узлам AST
			ast.Inspect(file, func(node ast.Node) bool {
				//pass.Reportf(file.Pos(),file.Name.Name)
				switch x := node.(type) {
				case *ast.ExprStmt: // выражение
					expr(x)

				case *ast.FuncDecl:
					if x.Name.String() == "main" {
						FuncMainsearchEnd(x)
					}

				case *ast.Ident:
					if x.Name == "os" {
						lastIdent = "os"
					} else {
						if lastIdent == "os" && x.Name == "Exit" {
							lastIdent = ""
							funcIdent(x)
						} else {
							lastIdent = ""
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
