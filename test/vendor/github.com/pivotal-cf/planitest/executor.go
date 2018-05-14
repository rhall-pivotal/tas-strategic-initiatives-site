package planitest

import (
	"bytes"
	"os"
	"os/exec"
)

type Executor struct {
	env []string
}

func NewExecutor() *Executor {
	return &Executor{}
}

func NewExecutorWithEnv(env []string) *Executor {
	return &Executor{
		env: env,
	}
}

func (e *Executor) Run(name string, args ...string) (string, string, error) {

	var errBuffer bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), e.env...)
	cmd.Stderr = &errBuffer
	output, err := cmd.Output()

	return string(output), errBuffer.String(), err
}
