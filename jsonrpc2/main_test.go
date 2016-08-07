package jsonrpc2

import (
	"testing"
	"github.com/kernelguy/gorpc"
	"github.com/kernelguy/gorpc/test"
)

func beginTest(name string) {
	test.Begin(name, 3)
}

func endTest() {
	test.End()
}

func TestMain(m *testing.M) {
	test.Init(m)
	gorpc.SetFactory(&Factory{})
	test.Run(m)
}
