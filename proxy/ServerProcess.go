package proxy

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type ServerProcess struct {
	running   bool
	Command   string
	Directory string
}

func (process *ServerProcess) Start() {
	if process.running {
		return
	}

	process.running = true

	log.Println("Starting Minecraft server...")
	InstallRedirect()

	commandParts := strings.Split(process.Command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = process.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Minecraft server shut down")
	UninstallRedirect()
}

func (process *ServerProcess) Stop() {
	log.Println("Stopping Minecraft server...")
}
