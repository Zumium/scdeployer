package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const DIRECT_SOURCE_LOADER_NAME = "DirectSourceLoader"
const DSL_SRC_FILE_NAME = "smartcontract.go"
const DSL_EXECUTABLE_FILE_NAME = "smartcontract"

type DirectSourceLoader struct {
	srcCode string
}

func (loader *DirectSourceLoader) ConfigByStrings(args []string) error {
	if len(args) != 1 {
		return errors.New("wrong arguments number")
	}
	return loader.SetParam(args[0])
}

func (loader *DirectSourceLoader) SetParam(param interface{}) error {
	var ok bool
	loader.srcCode, ok = param.(string)
	if !ok {
		return errors.New(fmt.Sprintf("%s reports: %s\n", DIRECT_SOURCE_LOADER_NAME, "parameter is not a string"))
	}
	if loader.srcCode == "" {
		return errors.New(fmt.Sprintf("%s reports: %s\n", DIRECT_SOURCE_LOADER_NAME, "illegal argument: empty source"))
	}
	return nil
}

func (loader *DirectSourceLoader) Load() (string, error) {
	f, err := os.Create(DSL_SRC_FILE_NAME)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, strings.NewReader(loader.srcCode)); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	compileCmd := exec.Command("go", "build", DSL_SRC_FILE_NAME)
	if err := compileCmd.Run(); err != nil {
		return "", err
	}
	//check whether target file has been generated successfully
	tgtFile, err := os.Open(DSL_EXECUTABLE_FILE_NAME)
	if err != nil {
		return "", err
	}
	if err := tgtFile.Close(); err != nil { // this should nerver happen
		return "", err
	}
	executableAbsPath, err := filepath.Abs(DSL_EXECUTABLE_FILE_NAME)
	if err != nil {
		return "", nil
	}
	return executableAbsPath, nil
}

func init() {
	loaders["direct_source_loader"] = func() ConfigByStringsLoader { return new(DirectSourceLoader) }
}
