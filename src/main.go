package main

import (
	"fmt"
	"twometer.dev/craftignite/minecraft"
	"twometer.dev/craftignite/proxy"
)

func main() {
	fmt.Println("CraftIgnite starting up")

	process := proxy.ServerProcess{
		FileName: "/app/minecraft.jar",
	}

	server := minecraft.Server{
		Motd:           "§6CraftIgnite Minecraft Proxy\n§7Server is currently sleeping",
		KickMessage:    "§l§6CraftIgnite\n\n§rThe server is currently starting.\nPlease try to reconnect in a minute.",
		TooltipMessage: "§aThis server is currently sleeping\n§rIt will automatically start once you join",
		MaxPlayerCount: 0,
		VersionName:    "1.0.0",
		ConnectCallback: func() {
			fmt.Println("A player tried to connect, starting Minecraft server...")
			process.Start()
		},
	}

	server.Start()
}
