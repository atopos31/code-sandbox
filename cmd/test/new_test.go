package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/atopos31/code-sandbox/internal/newcoder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
)

var Coder newcoder.Coder

func TestCoder(t *testing.T) {
	codetype := os.Args[len(os.Args)-1]
	var Coder newcoder.Coder

	switch codetype {
	case "cpp":
		code, _ := os.ReadFile("../testcode/test.cpp.txt")
		Coder = newcoder.NewCppCoder(string(code))
	case "c":
		code, _ := os.ReadFile("../testcode/test.c.txt")
		Coder = newcoder.NewCCoder(string(code))
	case "python":
		code, _ := os.ReadFile("../testcode/test.py.txt")
		Coder = newcoder.NewPyCoder(string(code))
	case "java":
		code, _ := os.ReadFile("../testcode/test.java.txt")
		Coder = newcoder.NewJavaCoder(string(code))
	case "go":
		code, _ := os.ReadFile("../testcode/test.go.txt")
		Coder = newcoder.NewGoCoder(string(code))
	default:
		t.Fatal("no codetype")
	}
	defer Coder.Clean()

	var SandBoxPool = sandbox.NewSandboxPool(100)
	sandbox, err := SandBoxPool.GetSandbox()
	if err != nil {
		t.Fatal(err)
	}
	defer SandBoxPool.ReleaseSandbox(sandbox)
	buildmeat, err := Coder.Build(sandbox)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("buildmeta: ", buildmeat)

	var stdins = []string{"100 3", "2 4", "4 6", "4 54", "4 9", "4 5", "4 6", "10 11", "22 22", "22 22"}
	var metas = []model.RunMeta{}

	var metac = make(chan model.RunMeta)
	for _, stdin := range stdins {
		go func() {
			sandboxr, _ := SandBoxPool.GetSandbox()
			defer SandBoxPool.ReleaseSandbox(sandbox)
			err := Coder.Run(sandboxr, stdin, 1, 10000000, metac)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	for i := 0; i < len(stdins); i++ {
		meta := <-metac
		metas = append(metas, meta)
	}

	for _, v := range metas {
		t.Log("meta: ", v)
	}
}
