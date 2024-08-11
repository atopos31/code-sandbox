package coder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
	"github.com/google/uuid"
)

type JavaCoder struct {
	sandbox      *sandbox.Sandbox
	basefielPath string
	binPath      string
	buildPath    string
	sidinPath    string
	MetaPath     string
	stderrPath   string
	stdoutPath   string
}

func NewJavaCoder() Coder {
	uuid := uuid.NewString()
	basePath := fmt.Sprintf("%s/%s", CodeStorageFolder, uuid)
	os.Mkdir(basePath, 0777)
	os.Chmod(basePath, 0777)
	return &JavaCoder{
		basefielPath: basePath,
		binPath:      fmt.Sprintf("Main"),
		buildPath:    fmt.Sprintf("%s/Main.java", basePath),
		sidinPath:    fmt.Sprintf("%s/stdin.txt", basePath),
		MetaPath:     fmt.Sprintf("%s/meta.txt", basePath),
		stderrPath:   fmt.Sprintf("%s/stderr.txt", basePath),
		stdoutPath:   fmt.Sprintf("%s/stdout.txt", basePath),
	}
}

func (c *JavaCoder) SetSandbox(sandbox *sandbox.Sandbox) {
	c.sandbox = sandbox
}

func (c *JavaCoder) Build(code string) (*model.CodeMETA, error) {
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
		"/usr/local/jdk17/bin/javac",
		c.buildPath,
	)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *JavaCoder) Run(MaxTime float64, MaxMem int, stdin string) (*model.CodeMETA, error) {
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
		"/usr/local/jdk17/bin/java",
		"-cp",
		c.basefielPath,
		c.binPath,
	)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	return model.NewCodeMETA(c.stderrPath, c.stdoutPath, c.MetaPath), nil
}

func (c *JavaCoder) Clean() {
	os.RemoveAll(c.basefielPath)
}
