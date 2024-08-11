package coder

import (
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
)

const CodeStorageFolder string = "/root/project/sandbox"

type Coder interface {
	SetSandbox(sandbox *sandbox.Sandbox)
	Build(code string) (*model.CodeMETA, error)
	Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error)
	Clean()
}
