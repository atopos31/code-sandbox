package coder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atopos31/code-sandbox/internal/model"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/google/uuid"
)

type GOCoder struct {
	sandbox      *sandbox.Sandbox
	basefielPath string
	binPath      string
	buildPath    string
	sidinPath    string
	MetaPath     string
	stderrPath   string
	stdoutPath   string
}

func NewGOCoder(sanbox *sandbox.Sandbox) *GOCoder {
	uuid := uuid.NewString()
	basePath := fmt.Sprintf("%s/%s", CodeStorageFolder, uuid)
	os.Mkdir(basePath, 0777)
	os.Chmod(basePath, 0777)
	return &GOCoder{
		sandbox:      sanbox,
		basefielPath: basePath,
		binPath:      fmt.Sprintf("%s/build", basePath),
		buildPath:    fmt.Sprintf("%s/build.go", basePath),
		sidinPath:    fmt.Sprintf("%s/stdin.txt", basePath),
		MetaPath:     fmt.Sprintf("%s/meta.txt", basePath),
		stderrPath:   fmt.Sprintf("%s/stderr.txt", basePath),
		stdoutPath:   fmt.Sprintf("%s/stdout.txt", basePath),
	}
}

func (g *GOCoder) Build(code string) (*model.CodeMETA, error) {
	os.WriteFile(g.buildPath, []byte(code), 0777)
	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", g.sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", g.basefielPath),
		"--dir=/root/.cache/go-build:rw",
		"--processes=100",
		"--fsize=5120",
		"--env=GOROOT",
		"--env=GOPATH",
		"--env=GO111MODULE=off",
		"--env=HOME",
		"--env=PATH",
		"--full-env",
		"--meta="+g.MetaPath,
		"--stdout="+g.stdoutPath,
		"--stderr="+g.stderrPath,
		"--wait",
		"--run",
		"--",
		"/usr/local/go/bin/go",
		"build",
		"-o",
		g.binPath,
		g.buildPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(g.stderrPath, g.stdoutPath, g.MetaPath), nil
}

func (g *GOCoder) Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error) {
	os.WriteFile(g.sidinPath, []byte(stdin), 0777)
	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", g.sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", g.basefielPath),
		"--processes=100",
		"--fsize=5120",
		"--env=PATH",
		"--meta="+g.MetaPath,
		fmt.Sprintf("--time=%f", MaxTime),
		fmt.Sprintf("--mem=%d", MaxMem),
		fmt.Sprintf("--stdin=%s", g.sidinPath),
		"--stdout="+g.stdoutPath,
		"--stderr="+g.stderrPath,
		"--wait",
		"--run",
		"--",
		g.binPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(g.stderrPath, g.stdoutPath, g.MetaPath), nil
}

func (g *GOCoder) Clean() {
	os.RemoveAll(g.basefielPath)
}

