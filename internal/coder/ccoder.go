package coder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atopos31/code-sandbox/internal/model"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/google/uuid"
)

type CCoder struct {
	sandbox      *sandbox.Sandbox
	basefielPath string
	binPath      string
	buildPath    string
	sidinPath    string
	MetaPath     string
	stderrPath   string
	stdoutPath   string
}

func NewCCoder(sanbox *sandbox.Sandbox) *CCoder {
	uuid := uuid.NewString()
	basePath := fmt.Sprintf("%s/%s", CodeStorageFolder, uuid)
	os.Mkdir(basePath, 0777)
	os.Chmod(basePath, 0777)
	return &CCoder{
		sandbox:      sanbox,
		basefielPath: basePath,
		binPath:      fmt.Sprintf("%s/build", basePath),
		buildPath:    fmt.Sprintf("%s/build.c", basePath),
		sidinPath:    fmt.Sprintf("%s/stdin.txt", basePath),
		MetaPath:     fmt.Sprintf("%s/meta.txt", basePath),
		stderrPath:   fmt.Sprintf("%s/stderr.txt", basePath),
		stdoutPath:   fmt.Sprintf("%s/stdout.txt", basePath),
	}
}

func (c *CCoder) Build(code string) (*model.CodeMETA, error) {
	os.WriteFile(c.buildPath, []byte(code), 0777)
	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", c.sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", c.basefielPath),
		"--processes=100",
		"--fsize=5120",
		"--env=PATH",
		"--meta="+c.MetaPath,
		"--stdout="+c.stdoutPath,
		"--stderr="+c.stderrPath,
		"--wait",
		"--run",
		"--",
		"/usr/bin/gcc",
		"-o",
		c.binPath,
		c.buildPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *CCoder) Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error) {
	os.WriteFile(c.sidinPath, []byte(stdin), 0777)
	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", c.sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", c.basefielPath),
		"--processes=100",
		"--fsize=5120",
		"--env=PATH",
		"--meta="+c.MetaPath,
		fmt.Sprintf("--time=%f", MaxTime),
		fmt.Sprintf("--mem=%d", MaxMem),
		fmt.Sprintf("--stdin=%s", c.sidinPath),
		"--stdout="+c.stdoutPath,
		"--stderr="+c.stderrPath,
		"--wait",
		"--run",
		"--",
		c.binPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *CCoder) Clean() {
	os.RemoveAll(c.basefielPath)
}
