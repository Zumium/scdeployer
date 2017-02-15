package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var loaderType string
var loaders map[string]ConfigByStringsLoaderGenerator = make(map[string]ConfigByStringsLoaderGenerator)

var rootCmd = &cobra.Command{
	Use: "scdeployer",
	Run: func(cmd *cobra.Command, args []string) {
		if loaderType == "" {
			printErrOnStderr(errors.New("must choose one loader"))
			return
		}
		generator, exists := loaders[loaderType]
		if !exists {
			printErrOnStderr(errors.New(fmt.Sprintf("no such loader: %s", loaderType)))
			return
		}
		loader := generator()

		if err := loader.ConfigByStrings(args); err != nil {
			printErrOnStderr(err)
			return
		}
		executablePath, err := loader.Load()
		if err != nil {
			printErrOnStderr(err)
			return
		}

		runCmd := exec.Command(executablePath)
		runCmd.Stdin = os.Stdin
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			printErrOnStderr(err)
			return
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		printErrOnStderr(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&loaderType, "loader", "l", "", "specify which loader to use")
}

func printErrOnStderr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
