package test

import (
	"fmt"
	"testing"

	"github.com/atopos31/code-sandbox/internal/newcoder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
)

var newCoder newcoder.NewCoder

func TestNewCpp(t *testing.T) {
	newCoder = newcoder.NewCppCoder("#include <iostream> \n using namespace std; int main() {int a,b;cin>>a>>b;cout<<a+b<<endl;}")
	defer newCoder.Clean()

	var SandBoxPool = sandbox.NewSandboxPool(100)
	sandbox, err := SandBoxPool.GetSandbox()
	if err != nil {
		t.Fatal(err)
	}
	defer SandBoxPool.ReleaseSandbox(sandbox)
	buildmeat, err := newCoder.Build(sandbox)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("buildmeta: ", buildmeat)

	var stdins = []string{"1 3", "2 4", "4 6","4 54","4 9","4 5","4 6","10 11"}
	var metas = []model.RunMeta{}
	var metac = make(chan model.RunMeta)

	for _, stdin := range stdins {
		go func() {
			sandboxr, _ := SandBoxPool.GetSandbox()
			defer SandBoxPool.ReleaseSandbox(sandbox)
			err := newCoder.Run(sandboxr, stdin, 1, 10240000, metac)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	for i := 0; i < len(stdins); i++ {
		meta := <-metac
		metas = append(metas, meta)
	}
	t.Log("metas: ", metas)
}
