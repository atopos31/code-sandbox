package main

import (
	"fmt"

	"github.com/atopos31/code-sandbox/internal/coder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
)

var SandBoxPool = sandbox.NewSandboxPool(10)

func main() {
	sanboxi, err := SandBoxPool.GetSandbox()
	if err != nil {
		panic(err)
	}
	defer SandBoxPool.ReleaseSandbox(sanboxi)

	coder := coder.NewCPPCoder(sanboxi)
	defer coder.Clean()

	buildMeta, err := coder.Build(`#include <iostream>
using namespace std;

int main() {
	int a,b;
	cin>>a>>b;
    cout<<a+b<<endl;
}
`)
	if err != nil {
		panic(err)
	}
	if buildMeta.Status != "" {
		fmt.Println("buildMeta", buildMeta)
		panic(buildMeta)
	}
	runMeta, err := coder.Run(1, 10000, " 111 9")
	if err != nil {
		panic(err)
	}
	fmt.Println("runMeta", runMeta)
}
