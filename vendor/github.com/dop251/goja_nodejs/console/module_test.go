package console

import (
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"testing"
)

func TestConsole(t *testing.T) {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	Enable(vm)

	if c := vm.Get("console"); c == nil {
		t.Fatal("console not found")
	}
}
