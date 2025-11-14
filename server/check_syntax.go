package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func main() {
	files := []string{
		"model/system/sys_user_action_log.go",
		"model/system/request/sys_user_action_log.go",
		"model/system/response/sys_user_action_log.go",
		"service/system/sys_user_action_log.go",
		"api/v1/system/sys_user_action_log.go",
		"router/system/sys_user_action_log.go",
	}

	hasError := false
	for _, file := range files {
		fullPath := filepath.Join(".", file)
		if err := checkFile(fullPath); err != nil {
			fmt.Printf("❌ %s: %v\n", file, err)
			hasError = true
		} else {
			fmt.Printf("✅ %s: OK\n", file)
		}
	}

	if hasError {
		os.Exit(1)
	}
	fmt.Println("\n所有文件语法检查通过！")
}

func checkFile(filename string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return err
	}
	return nil
}
