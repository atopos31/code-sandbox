package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Sandbox struct {
	boxID int
}

func NewSandbox() *Sandbox {
	i := 0
	bcmd := exec.Command("isolate", "--init", fmt.Sprintf("-b %v", i))
	op, err := bcmd.Output()
	if err != nil {
		fmt.Printf("failed to create sandbox %v, output: %v err: %v", i, string(op), err)
	}

	return &Sandbox{
		boxID: i,
	}
}

func (s *Sandbox) RunCPP(filepath string) {
	// Define the path for the meta file to store execution information
	metaFilePath := fmt.Sprintf("/root/project/sandbox/meta_%d.txt", s.boxID)
	stdoutFilePath := fmt.Sprintf("/root/project/sandbox/stdout_%d.txt", s.boxID)
	stderr := fmt.Sprintf("/root/project/sandbox/stderr_%d.txt",s.boxID)
	stdin := fmt.Sprintf("/root/project/sandbox/stdin_%d.txt",s.boxID)


	// Compile the C++ program
	bcmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", s.boxID),
		"--dir=/root/project/sandbox:rw",
		"--processes=10",
		"--fsize=5120",
		"--env=PATH",
		"--meta="+metaFilePath,
		"--stdout="+stdoutFilePath,
		"--stderr="+stderr,
		"--wait",
		"--run",
		"--",
		"/usr/bin/g++",
		"-o",
		"test",
		filepath,
	)
	bcmd.Run()

	//Run the compiled program
	rcmd := exec.Command("isolate",
		fmt.Sprintf("--box-id=%d", s.boxID),
		"--dir=/root/project/sandbox:rw",
		"--meta="+metaFilePath,
		"--env=PATH",
		"--stdout="+stdoutFilePath,
		"--stdin="+ stdin,
		"--stderr="+stderr,
		"--wait",
		"--run",
		"--",
		"./test",
	)
	rcmd.Run()

	data,err  := os.Open(stdoutFilePath)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	var lines []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查扫描器是否遇到错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 拼接所有行，并去除最后一个换行符
	result := strings.Join(lines, "\n")
	fmt.Println(result=="36\n36")

	// Read and display the meta information
	metaData, err := os.ReadFile(metaFilePath)
	if err != nil {
		fmt.Printf("Failed to read meta information: %v\n", err)
		return
	}
	fmt.Printf("Execution meta data:\n%v\n", string(metaData))
}
