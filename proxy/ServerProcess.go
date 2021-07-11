package proxy

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ServerProcess struct {
	running          bool
	stdin            io.WriteCloser
	Command          string
	Directory        string
	ShutdownCallback func()
	StartupCallback  func()
}

func (process *ServerProcess) Start() {
	if process.running {
		return
	}

	process.running = true

	log.Println("Starting Minecraft server...")
	log.Println("Server command line: " + process.Command)

	// Start the server
	commandParts := strings.Split(process.Command, " ")
	cmd := exec.Command(commandParts[0], commandParts[1:]...)
	cmd.Dir = process.Directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}

	// Start process
	process.stdin = stdin
	err = cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	// Wait for the server to reply
	for {
		_, err := GetServerStatus()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	process.StartupCallback()

	// STDIN Passthrough
	go func() {
		_, err := io.Copy(process.stdin, os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Wait for server to shut down
	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Minecraft server shut down")
	process.running = false
	process.ShutdownCallback()
}

func (process *ServerProcess) Stop() {
	log.Println("Stopping Minecraft server...")
	_, err := process.stdin.Write(([]byte)("stop\r\n"))
	if err != nil {
		log.Fatalln(err)
	}
}
