package main

import (
	"os"
	"os/exec"
)

type commander interface {
	exec(string, ...string) ([]byte, error)
	execWithOutput(command string, args ...string) ([]byte, error)
}

type realCommander struct {
}

func (c realCommander) exec(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return []byte{}, err
}

func (c realCommander) execWithOutput(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	return out, err
}
