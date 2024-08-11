package coder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atopos31/code-sandbox/internal/model"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/google/uuid"
)

type PythonCoder struct {
	sandbox      *sandbox.Sandbox
	basefielPath string
	binPath      string
	buildPath    string
	sidinPath    string
	MetaPath     string
	stderrPath   string
	stdoutPath   string
}

func NewPythonCoder(sanbox *sandbox.Sandbox) *PythonCoder {
	uuid := uuid.NewString()
	basePath := fmt.Sprintf("%s/%s", CodeStorageFolder, uuid)
	os.Mkdir(basePath, 0777)
	os.Chmod(basePath, 0777)
	return &PythonCoder{
		sandbox:      sanbox,
		basefielPath: basePath,
		binPath:      fmt.Sprintf("%s/__pycache__/build.cpython-38.pyc", basePath),
		buildPath:    fmt.Sprintf("%s/build.py", basePath),
		sidinPath:    fmt.Sprintf("%s/stdin.txt", basePath),
		MetaPath:     fmt.Sprintf("%s/meta.txt", basePath),
		stderrPath:   fmt.Sprintf("%s/stderr.txt", basePath),
		stdoutPath:   fmt.Sprintf("%s/stdout.txt", basePath),
	}
}

func (c *PythonCoder) Build(code string) (*model.CodeMETA, error) {
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
		"/usr/bin/python3",
		"-m",
		"compileall",
		c.buildPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *PythonCoder) Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error) {
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
		"/usr/bin/python3",
		c.binPath,
	)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *PythonCoder) Clean() {
	os.RemoveAll(c.basefielPath)
}
