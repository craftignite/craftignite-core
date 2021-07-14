package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"twometer.dev/craftignite/minecraft"
	"twometer.dev/craftignite/proxy"
)

func main() {
	log.Println("CraftIgnite starting up")

	timeoutStr := os.Getenv("CRAFTIGNITE_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)

	if err != nil {
		log.Fatalln(err)
	}

	var process proxy.ServerProcess

	watchdog := proxy.Watchdog{
		Timeout: int64(timeout),
		ShutdownCallback: func() {
			process.Stop()
		},
	}
	go watchdog.Start()

	server := minecraft.Server{
		Motd:           os.Getenv("CRAFTIGNITE_MOTD"),
		KickMessage:    os.Getenv("CRAFTIGNITE_KICK_MESSAGE"),
		TooltipMessage: os.Getenv("CRAFTIGNITE_TOOLTIP_MESSAGE"),
		HostAddress:    ":" + os.Getenv("SERVER_PORT"),
		MaxPlayerCount: 0,
		VersionName:    "1.0.0",
		ConnectCallback: func() {
			watchdog.Reset()
			go process.Start()
		},
	}

	process = proxy.ServerProcess{
		Command:   strings.Join(os.Args[1:], " "),
		Directory: ".",
		StartupCallback: func() {
			server.ProxyMode = true
			watchdog.Reset()
		},
		ShutdownCallback: func() {
			server.ProxyMode = false
			if !watchdog.HasShutdown {
				log.Println("Server stopped without watchdog, interpreting as full stop")
				os.Exit(0)
			}
		},
	}

	server.Start()
}
