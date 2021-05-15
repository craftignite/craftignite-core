package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ServerProcess struct {
	starting  bool
	Command   string
	Directory string
}

func (process *ServerProcess) Start() {
	if process.starting {
		return
	}

	process.starting = true

	fmt.Println("Starting Minecraft server...")

	commandParts := strings.Split(process.Command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = process.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
