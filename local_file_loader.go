package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

const LOCAL_FILE_LOADER_NAME = "LocalFileLoader"

type LocalFileLoader struct {
	srcFilePath string
}

func (loader *LocalFileLoader) ConfigByStrings(args []string) error {
	if len(args) != 1 {
		return errors.New("wrong arguments number")
	}
	return loader.SetParam(args[0])
}

func (loader *LocalFileLoader) SetParam(param interface{}) error {
	var ok bool
	loader.srcFilePath, ok = param.(string)
	if !ok {
		return errors.New(fmt.Sprintf("%s reports: %s\n", LOCAL_FILE_LOADER_NAME, "parameter is not a string"))
	}
	if loader.srcFilePath == "" {
		return errors.New(fmt.Sprintf("%s reports: %s\n", LOCAL_FILE_LOADER_NAME, "illegal argument: null file path"))
	}
	return nil
}

func (loader *LocalFileLoader) Load() (string, error) {
	compileCmd := exec.Command("go", "build", "-o", "smartcontract", loader.srcFilePath)
	if err := compileCmd.Run(); err != nil {
		return "", err
	}
	executableAbsPath, err := filepath.Abs("smartcontract")
	if err != nil {
		return "", nil
	}
	return executableAbsPath, nil
}

func init() {
	loaders["local_file_loader"] = func() ConfigByStringsLoader { return new(LocalFileLoader) }
}
