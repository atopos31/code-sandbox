package coder

import (
	"github.com/atopos31/code-sandbox/internal/model"
)

const CodeStorageFolder string = "/root/project/sandbox"

type Coder interface {
	Build(code string) (*model.CodeMETA, error)
	Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error)
}

