package proxy

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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

	// Start the server
	commandParts := strings.Split(process.Command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = process.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	// Wait for the server to reply
	for {
		_, err := GetServerStatus("127.0.0.1")
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	InstallRedirect()
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
