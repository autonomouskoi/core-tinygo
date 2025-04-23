//go:build mage

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"

	svc "github.com/autonomouskoi/core-tinygo/svc/build"
	"github.com/autonomouskoi/mageutil"
)

var (
	baseDir     string
	testWASMDir string
	wasmDir     string
)

var Default = Protos

func init() {
	var err error
	baseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	baseDir = filepath.Join(baseDir, "..")
	testWASMDir = filepath.Join(baseDir, "test_wasm")
	wasmDir = filepath.Join(baseDir, "wasm_out")
}

func Clean() error {
	return sh.Rm(wasmDir)
}

func Protos() {
	mg.Deps(
		TinyGoProtos,
		svc.TinyGoProtos,
	)
}

func wasmOutDir() error {
	return mageutil.Mkdir(wasmDir)
}

/*
func TestWASM() error {
	mg.Deps(Dev, wasmOutDir)
	return mageutil.WasmTinygoDir(filepath.Join(wasmDir, "test.wasm"), testWASMDir)
}
*/

func TinyGoProtos() error {
	akdir := filepath.Join(baseDir, "..")
	srcpath := filepath.Join(akdir, "akcore", "bus", "bus.proto")
	outfile := "bus.pb.go"
	outpath := filepath.Join(baseDir, outfile)

	newer, err := target.Path(outpath, srcpath)
	if err != nil {
		return fmt.Errorf("testing %s vs %s: %w", srcpath, outpath, err)
	}
	if !newer {
		return nil
	}

	pcggl, err := exec.LookPath("protoc-gen-go-lite")
	if err != nil {
		return fmt.Errorf("finding protoc plugin: %w", err)
	}
	err = sh.Run("protoc",
		"--plugin", pcggl,
		"--go-lite_opt", "features=marshal+unmarshal+size+equal+clone",
		"-I", filepath.Join(akdir, "akcore"),
		"--go-lite_out", baseDir,
		srcpath,
	)
	if err != nil {
		return err
	}
	genfile := filepath.Join(baseDir, "github.com", "autonomouskoi", "akcore", "bus", outfile)
	genfh, err := os.Open(genfile)
	if err != nil {
		return fmt.Errorf("opening %s: %w", genfile, err)
	}
	defer genfh.Close()
	outfh, err := os.Create(outpath)
	if err != nil {
		return fmt.Errorf("creating %s: %w", outpath, err)
	}
	defer outfh.Close()
	scanner := bufio.NewScanner(genfh)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "package bus" {
			text = "package core"
		}
		fmt.Fprintln(outfh, text)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning: %w", err)
	}
	if err := outfh.Sync(); err != nil {
		return fmt.Errorf("syncing %s: %w", outpath, err)
	}
	if err := sh.Rm(filepath.Join(baseDir, "github.com")); err != nil {
		return fmt.Errorf("removing: %w", err)
	}
	return nil
}
