package main

import (
	"github.com/atopos31/code-sandbox/internal/app"
	"github.com/atopos31/code-sandbox/internal/sandbox"
)

func main() {
    app := app.New(sandbox.NewSandboxPool(100))
    app.Run()
}