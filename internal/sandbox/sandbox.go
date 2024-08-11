package sandbox

import (
	"fmt"
	"os/exec"
)

type Sandbox struct {
	ID int
}

func newSandbox(id int) (*Sandbox, error) {
	cmd := exec.Command("isolate", "--init", fmt.Sprintf("-b %v", id))
	op, err := cmd.Output()
	if err != nil {
		fmt.Printf("failed to create sandbox %v, output: %v err: %v", id, string(op), err)
		return nil, err
	}
	return &Sandbox{ID: id}, nil
}
