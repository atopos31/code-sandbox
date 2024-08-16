package newcoder

import (
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
)

const codeStorageFolder string = "/sandbox/running"

type NewCoderFunc func(cppCode string) Coder

type Coder interface {
	Build(sandbox *sandbox.Sandbox) (*model.BuildMeta, error)
	Run(sandbox *sandbox.Sandbox, stdin string, MaxTime float64, MaxMem int, meta chan<- model.RunMeta) error
	Clean()
}
