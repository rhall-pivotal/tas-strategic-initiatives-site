package planitest

import (
	"bytes"
	"os/exec"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Run(name string, args ...string) (string, string, error) {

	var errBuffer bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stderr = &errBuffer
	output, err := cmd.Output()

	return string(output), errBuffer.String(), err
}
