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

	meta, err := coder.Build(`#include <iostream>
using namespace std

int main() {
    cout<<"test"<<endl;
}
`)
	if err != nil {
		panic(err)
	}
	fmt.Println(meta)
}
