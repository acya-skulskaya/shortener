package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const (
	MsgUsingPanic    = "Using panic function is discouraged"
	MsgUsingLogFatal = "Using log.Fatal function outside of main function of main package is discouraged"
	MsgUsingOsExit   = "Using os.Exit function outside of main function of main package is discouraged"
)

var Analyzer = &analysis.Analyzer{
	Name: "pfeLinter",
	Doc:  "checks for panic() function calls in project code, and for log.Fatal() and/or os.Exit() function calls outside of main function of main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch call := node.(type) {
			case *ast.CallExpr:
				checkForPanic(pass, call)
				checkForLogFatal(pass, call, file)
				checkForOsExit(pass, call, file)
			}

			return true
		})
	}
	return nil, nil
}

func checkForPanic(pass *analysis.Pass, call *ast.CallExpr) {
	if fun, ok := call.Fun.(*ast.Ident); ok && fun.Name == "panic" {
		pass.Reportf(call.Fun.Pos(), MsgUsingPanic)
	}
}

func checkForLogFatal(pass *analysis.Pass, call *ast.CallExpr, file *ast.File) {
	if checkPkgAndFuncName(call, "log", "Fatal") {
		if !checkIsMain(pass, call, file) {
			pass.Reportf(call.Fun.Pos(), MsgUsingLogFatal)
		}
	}
}
func checkForOsExit(pass *analysis.Pass, call *ast.CallExpr, file *ast.File) {
	if checkPkgAndFuncName(call, "os", "Exit") {
		if !checkIsMain(pass, call, file) {
			pass.Reportf(call.Fun.Pos(), MsgUsingOsExit)
		}
	}
}

func checkPkgAndFuncName(call *ast.CallExpr, pkgName string, funcName string) bool {
	se, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	i, ok := se.X.(*ast.Ident)
	if !ok {
		return false
	}

	if i.Name == pkgName && se.Sel.Name == funcName {
		return true
	}

	return false
}

func checkIsMain(pass *analysis.Pass, call *ast.CallExpr, f *ast.File) bool {
	if pass.Pkg == nil {
		return false
	}

	if pass.Pkg.Name() != "main" {
		return false
	}

	if f.Name.Name != "main" {
		return true
	}

	pos := call.Pos()
	file := pass.Fset.File(pos)
	line := file.Line(pos)

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "main" {
			continue
		}

		startLine := pass.Fset.Position(fn.Body.Pos()).Line
		endLine := pass.Fset.Position(fn.Body.End()).Line

		if file.Line(pos) >= startLine && line <= endLine {
			return true
		}
	}

	return false
}
