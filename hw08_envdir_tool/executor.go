package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	command.Env = []string{}

	for name, envVar := range env {
		if !envVar.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", name, envVar.Value))
		}
	}

	// rc := command.Run()
	//
	//	if rc == nil {
	//		return 0
	//	}
	//
	// return rc.(exec.ExitError)

	if err := command.Start(); err != nil {
		log.Fatal(err)
	}
	if err := command.Wait(); err != nil {
		log.Fatal(err)
	}
	rc := command.ProcessState.ExitCode()
	return rc
}
