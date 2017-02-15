package main

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

const NET_REPO_LOADER_NAME = "NetRepoLoader"

type NetRepoLoader struct {
	repo string
}

func (loader *NetRepoLoader) ConfigByStrings(args []string) error {
	if len(args) != 1 {
		return errors.New("wrong arguments number")
	}
	return loader.SetParam(args[0])
}

func (loader *NetRepoLoader) SetParam(param interface{}) error {
	var ok bool
	loader.repo, ok = param.(string)
	if !ok {
		return errors.New(fmt.Sprintf("%s reports: %s", NET_REPO_LOADER_NAME, "parameter is not a string"))
	}
	if loader.repo == "" {
		return errors.New(fmt.Sprintf("%s reports: %s", NET_REPO_LOADER_NAME, "illegal argument: null repo string"))
	}
	return nil
}

func (loader *NetRepoLoader) Load() (string, error) {
	downloadCmd := exec.Command("go", "get", loader.repo)
	if err := downloadCmd.Run(); err != nil {
		return "", err
	}
	compileCmd := exec.Command("go", "build", "-o", "smartcontract", loader.repo)
	if err := compileCmd.Run(); err != nil {
		return "", err
	}
	executableAbsPath, err := filepath.Abs("smartcontract")
	if err != nil {
		return "", err
	}
	return executableAbsPath, nil
}

func init() {
	loaders["net_repo_loader"] = func() ConfigByStringsLoader { return new(NetRepoLoader) }
}
