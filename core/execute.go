package core

import (
	"bytes"
	"os/exec"
	"io"
	"os"
	"log"
	"strings"
)

func Execute(command string) {
	ExecuteIn(command, "")
}

func ExecuteIn(command string, directory string) {
	args := strings.Split(command, " ")

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = directory

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
}
