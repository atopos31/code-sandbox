package test

import (
	"os"
	"testing"

	"github.com/atopos31/code-sandbox/internal/coder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
)

func TestGO(t *testing.T) {
	var SandBoxPool = sandbox.NewSandboxPool(10)
	sandbox, err := SandBoxPool.GetSandbox()
	if err != nil {
		t.Fatal(err)
	}
	defer SandBoxPool.ReleaseSandbox(sandbox)

	coder := coder.NewGOCoder()
	coder.SetSandbox(sandbox)
	defer coder.Clean()
	code, err := os.ReadFile("../testcode/test.go.txt")
	if err != nil {
		t.Fatal(err)
	}
	meta, err := coder.Build(string(code))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Status != "" {
		t.Fatal(meta)
	}
	meta, err = coder.Run(1, 1000000, "3 4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(meta)
}
