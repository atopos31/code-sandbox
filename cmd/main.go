package main

import (
	"os"

	"github.com/atopos31/code-sandbox/internal/app"
	"github.com/atopos31/code-sandbox/internal/sandbox"
)

func main() {
	port := os.Getenv("SERVICE_PORT")
	sandBoxPool := sandbox.NewSandboxPool(100)
	app := app.New(sandBoxPool)
	app.Run(port)
}
