package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/moznion/gowrtr/generator"
)

var fset = token.NewFileSet()

func main() {

	filepath := os.Args[1]
	filenameWithExt := path.Base(filepath)
	filename := filenameWithExt[0 : len(filenameWithExt)-3]
	f, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// register Func
	registerRoutesFunc := generateNewFunc("RegisterRoutes", "e", "*echo.Echo")
	registerAuthRoutesFunc := generateNewFunc("RegisterAuthRoutes", "e", "*echo.Group")

	type Group struct {
		Name   string
		Fields []*ast.Field
	}

	var groups []*Group
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			groupName := strings.Replace(typeSpec.Name.Name, "Routes", "", -1)
			fields := structType.Fields.List
			group := new(Group)
			group.Name = groupName
			group.Fields = fields
			groups = append(groups, group)
		}
	}

	var registererFuncs []generator.Statement
	for _, v := range groups {
		// get info from apis
		//registererFuncs := make([]generator.Statement, len(v.Fields))
		var httpMethods []string
		var apiPaths []string
		var isAuths []bool
		var funcNames []string
		for _, field := range v.Fields {
			tagValue := field.Tag.Value[1 : len(field.Tag.Value)-1]
			tag := reflect.StructTag(tagValue)

			httpMethod := tag.Get("method")
			if httpMethod == "" {
				panic("method mustn't be empty")
			}
			apiPath := tag.Get("path")
			if apiPath == "" {
				panic("path mustn't be empty")
			}
			funcName := field.Names[0].Name

			auth := tag.Get("auth")
			isAuth := auth == "true"
			httpMethods = append(httpMethods, httpMethod)
			apiPaths = append(apiPaths, apiPath)
			isAuths = append(isAuths, isAuth)
			funcNames = append(funcNames, funcName)
		}

		// generate func for apis
		for i, funcName := range funcNames {
			anon, _ := generatorNewAnonymousFunc(isAuths[i], funcName)
			anon = strings.TrimSpace(anon)

			var funcParameter1 *generator.FuncParameter
			if isAuths[i] {
				funcParameter1 = generator.NewFuncParameter("e", "*echo.Group")
			} else {
				funcParameter1 = generator.NewFuncParameter("e", "*echo.Echo")
			}
			funcParameter2 := generator.NewFuncParameter("inter", "*controller."+v.Name)

			registererFunc := generator.NewFunc(
				nil,
				generator.NewFuncSignature(funcName).AddParameters(funcParameter1, funcParameter2),
				generator.NewRawStatementf("e.%s", httpMethods[i]).WithNewline(false),
				generator.NewFuncInvocation(fmt.Sprintf(`"%s"`, apiPaths[i]), anon),
			)
			registererFuncs = append(registererFuncs, registererFunc)

			rawStatementf := generator.NewRawStatementf("%s(e, &controller.%s{})", funcName, v.Name)
			if isAuths[i] {
				registerAuthRoutesFunc = registerAuthRoutesFunc.AddStatements(rawStatementf)
			} else {
				registerRoutesFunc = registerRoutesFunc.AddStatements(rawStatementf)
			}
		}
	}

	src, err := generator.NewRoot(
		generator.NewComment(" This file was auto-generated."),
		generator.NewComment(" DO NOT EDIT MANUALLY!!!"),
		generator.NewPackage("api"),
		generator.NewImport("github.com/midnight-trigger/todo/api/controller"),
		registerRoutesFunc,
		registerAuthRoutesFunc,
	).AddStatements(registererFuncs...).Gofmt("-s").Goimports().Generate(0)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s_gen.go", path.Dir(filepath), filename), []byte(src), 0644)
	if err != nil {
		panic(err)
	}
}

func generateNewFunc(funcName string, parametersName string, parametersType string) *generator.Func {
	return generator.NewFunc(
		nil,
		generator.NewFuncSignature(funcName).AddParameters(generator.NewFuncParameter(parametersName, parametersType)),
	)
}

func generatorNewAnonymousFunc(isAuth bool, funcName string) (string, error) {
	funcSignature := generator.NewAnonymousFuncSignature().AddParameters(generator.NewFuncParameter("c", "echo.Context")).ReturnTypes("error")

	var jwtRawStatement *generator.RawStatement
	var interRawStatement *generator.RawStatement
	if isAuth {
		jwtRawStatement = generator.NewRawStatement("claims,r := jwt.GetJWTClaims(c); if claims == nil { return c.JSON(r.Code, r) }")
		interRawStatement = generator.NewRawStatementf("res := inter.%s(c, claims)", funcName)
	} else {
		jwtRawStatement = generator.NewRawStatement("").WithNewline(false)
		interRawStatement = generator.NewRawStatementf("res := inter.%s(c)", funcName)
	}

	returnStatement := generator.NewReturnStatement("c.JSON(res.Meta.Code, res)")
	return generator.NewAnonymousFunc(false, funcSignature, jwtRawStatement, interRawStatement, returnStatement).Generate(0)
}
