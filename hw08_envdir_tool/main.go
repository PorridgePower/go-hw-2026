package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// $ go-envdir /path/to/env/dir command arg1 arg2
	args := pflag.Args()
	if len(args) < 2 {
		fmt.Println("Not enough args")
		os.Exit(1)
	}
	envDir := args[0]
	// cmdName := args[1]
	// cmdArgs := args[2:]

	//====================
	// cmd := exec.Command(cmdName, cmdArgs...)
	// RunCmd(cmd, envDir)
	// cmdString := args[1:]
	// RunCmd(cmdString, envDir)
	//====================
	// dir, err := os.ReadDir(envDir)
	// if err != nil {
	// 	fmt.Println("Cannot read dir")
	// 	os.Exit(1)
	// }
	// env := make(Environment)
	// for _, item := range dir {
	// 	if item.IsDir() {
	// 		continue
	// 	}
	// 	fname := item.Name()
	// 	if strings.Contains(fname, "=") {
	// 		continue
	// 	}
	// 	data, err := os.ReadFile(filepath.Join(envDir, fname))
	// 	if err != nil {
	// 		fmt.Println("Cannot read file")
	// 		os.Exit(1)
	// 	}
	// 	binaryLine := bytes.SplitN(data, []byte("\n"), 2)[0]
	// 	binaryLine = bytes.Replace(binaryLine, []byte("\x00"), []byte("\n"), -1)
	// 	varValue := strings.TrimRight(string(binaryLine), " \t")
	// 	env[fname] = EnvValue{Value: varValue, NeedRemove: false}
	// }
	env, err := ReadDir(envDir)
	if err != nil {
		os.Exit(1)
	}
	cmd := args[1:]
	RunCmd(cmd, env)
}
