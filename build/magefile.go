package bus

import (
	"path/filepath"

	build "github.com/autonomouskoi/akcore/build/common"
	"github.com/autonomouskoi/mageutil"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	baseDir     string
	testWASMDir string
	wasmDir     string
)

func init() {
	baseDir = filepath.Join(build.BaseDir, "bus", "tinygo")
	testWASMDir = filepath.Join(baseDir, "test_wasm")
	wasmDir = filepath.Join(baseDir, "wasm_out")
}

func Clean() error {
	return sh.Rm(wasmDir)
}

func Dev() {}

func wasmOutDir() error {
	return mageutil.Mkdir(wasmDir)
}

func TestWASM() error {
	mg.Deps(Dev, wasmOutDir)
	return mageutil.WasmTinygoDir(filepath.Join(wasmDir, "test.wasm"), testWASMDir)
}
