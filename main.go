package main

import (
	"log"
	"twometer.dev/craftignite/minecraft"
	"twometer.dev/craftignite/proxy"
)

func main() {
	log.Println("CraftIgnite starting up")

	process := proxy.ServerProcess{
		Command:   "java -jar server.jar",
		Directory: ".testserver/",
	}

	watchdog := proxy.Watchdog{
		Timeout: 30,
		ShutdownCallback: func() {
			process.Stop()
		},
	}
	go watchdog.Start()

	server := minecraft.Server{
		Motd:           "§6CraftIgnite Minecraft Proxy\n§7Server is currently sleeping",
		KickMessage:    "§l§6CraftIgnite\n\n§rThe server is currently starting.\nPlease try to reconnect in a minute.",
		TooltipMessage: "§aServer will automatically start once you join",
		MaxPlayerCount: 0,
		VersionName:    "1.0.0",
		ConnectCallback: func() {
			watchdog.Reset()
			go process.Start()
		},
	}
	server.Start()
}
