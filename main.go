package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethpollack/dockerbox/api"
	"github.com/sethpollack/dockerbox/io"
)

func main() {
	rootDir := io.GetEnv("DOCKERBOX_ROOT_DIR", "$HOME/.dockerbox")
	installDir := io.GetEnv("DOCKERBOX_INSTALL_DIR", rootDir+"/bin")
	separator := io.GetEnv("DOCKERBOX_SEPARATOR", "--")

	entrypoint := filepath.Base(os.Args[0])
	args := os.Args[1:]
	exe, _ := os.Executable()

	err := api.New(rootDir, installDir, separator, exe, entrypoint, args).Run()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
