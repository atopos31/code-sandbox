package newcoder

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
	"github.com/google/uuid"
)

type GoCoder struct {
	basePath      string // coder基本路径 {codeStorageFolder}/{uuid}
	codePath      string // 源代码路径
	binPath       string // 编译后二进制文件路径
	buildMetaPath string // 编译元数据路径
	buildErrPath  string // 编译错误信息路径
	runsBasePath  string // 运行基本路径
}

func NewGoCoder(cppCode string) Coder {
	basePath := fmt.Sprintf("%s/%s", codeStorageFolder, uuid.NewString())
	os.Mkdir(basePath, 0777)
	os.Chmod(basePath, 0777)

	codePath := fmt.Sprintf("%s/code.go", basePath)
	os.WriteFile(codePath, []byte(cppCode), 0777)

	binPath := fmt.Sprintf("%s/bin", basePath)
	buildMetaPath := fmt.Sprintf("%s/buildMeta.txt", basePath)
	buildErrPath := fmt.Sprintf("%s/buildErr.txt", basePath)

	runsPath := fmt.Sprintf("%s/runs", basePath)
	os.Mkdir(runsPath, 0777)
	os.Chmod(runsPath, 0777)

	return &GoCoder{
		basePath:      basePath,
		codePath:      codePath,
		binPath:       binPath,
		buildMetaPath: buildMetaPath,
		buildErrPath:  buildErrPath,
		runsBasePath:  runsPath,
	}
}

func (c *GoCoder) Build(sandbox *sandbox.Sandbox) (*model.BuildMeta, error) {
	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", c.basePath),
		"--dir=/root/.cache/go-build:rw",
		"--processes=100",
		"--fsize=5120",
		"--env=GOROOT",
		"--env=GOPATH",
		"--env=GO111MODULE=off",
		"--env=HOME",
		"--env=PATH",
		"--meta="+c.buildMetaPath,
		"--stderr="+c.buildErrPath,
		"--wait",
		"--run",
		"--",
		"/usr/local/go/bin/go",
		"build",
		"-o",
		c.binPath,
		c.codePath,
	)

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	buildMeta := model.NewBuildMeta(c.buildErrPath, c.buildMetaPath)

	return buildMeta, nil
}

// 并行运行程序
func (c *GoCoder) Run(sandbox *sandbox.Sandbox, stdin string, MaxTime float64, MaxMem int, meta chan<- model.RunMeta) error {
	runPath := fmt.Sprintf("%s/%s", c.runsBasePath, uuid.NewString())
	os.Mkdir(runPath, 0777)
	os.Chmod(runPath, 0777)
	stdinPath := fmt.Sprintf("%s/stdin.txt", runPath)
	os.WriteFile(stdinPath, []byte(stdin), 0777)

	stdoutPath := fmt.Sprintf("%s/stdout.txt", runPath)
	stderrPath := fmt.Sprintf("%s/stderr.txt", runPath)
	metaPath := fmt.Sprintf("%s/meta.txt", runPath)

	cmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", sandbox.ID),
		fmt.Sprintf("--dir=%s:rw", c.basePath),
		"--processes=100",
		"--fsize=5120",
		"--env=PATH",
		"--meta="+metaPath,
		fmt.Sprintf("--time=%f", MaxTime),
		fmt.Sprintf("--mem=%d", MaxMem),
		fmt.Sprintf("--stdin=%s", stdinPath),
		"--stdout="+stdoutPath,
		"--stderr="+stderrPath,
		"--wait",
		"--run",
		"--",
		c.binPath,
	)

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	metadata := model.NewRunMeta(stderrPath, stdoutPath, metaPath)
	metadata.Stdin = stdin
	meta <- *metadata
	return nil
}

func (c *GoCoder) Clean() {
	os.RemoveAll(c.basePath)
}
