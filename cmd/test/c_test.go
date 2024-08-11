package test

import (
	"os"
	"testing"

	"github.com/atopos31/code-sandbox/internal/coder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
)

func TestC(t *testing.T) {
	var SandBoxPool = sandbox.NewSandboxPool(10)
	sandbox, err := SandBoxPool.GetSandbox()
	if err != nil {
		t.Fatal(err)
	}
	defer SandBoxPool.ReleaseSandbox(sandbox)

	coder := coder.NewCCoder(sandbox)
	defer coder.Clean()
	code, err := os.ReadFile("../testcode/test.c.txt")
	if err != nil {
		t.Fatal(err)
	}
	meta, err := coder.Build(string(code))
	if err != nil {
		t.Fatal(err)
	}
	if meta.Status != "" {
		t.Log(meta)
		return
	}

	meta, err = coder.Run(1, 10000, "3 4")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(meta)

}
