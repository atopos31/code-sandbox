package test

import (
	"os"
	"testing"
)

func TestMakeDir(t *testing.T) {
	os.Mkdir("/root/project/sandbox/test", 0777)
	os.Chmod("/root/project/sandbox/test", 0777)
}
