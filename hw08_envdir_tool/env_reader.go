package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envDir, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Cannot read dir")
		os.Exit(1)
	}
	env := make(Environment)
	for _, item := range envDir {
		if item.IsDir() {
			continue
		}
		fname := item.Name()
		if strings.Contains(fname, "=") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, fname))
		if err != nil {
			fmt.Println("Cannot read file")
			os.Exit(1)
		}
		binaryLine := bytes.SplitN(data, []byte("\n"), 2)[0]
		binaryLine = bytes.Replace(binaryLine, []byte("\x00"), []byte("\n"), -1)
		varValue := strings.TrimRight(string(binaryLine), " \t")
		env[fname] = EnvValue{Value: varValue, NeedRemove: false}

		// TODO: envdir удаляет переменную среды с именем S, если таковая существует
	}
	return nil, nil
}
