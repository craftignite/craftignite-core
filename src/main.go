package main

import (
	"fmt"
	"twometer.dev/craftignite/minecraft"
)

func main() {
	fmt.Println("CraftIgnite starting up")
	server := minecraft.Server{
		Motd:           "§6CraftIgnite Minecraft Proxy\n§7Server is currently sleeping",
		KickMessage:    "§l§6CraftIgnite\n\n§rThe server is currently starting.\nPlease try to reconnect in a minute.",
		TooltipMessage: "§aThis server is currently sleeping\n§rIt will automatically start once you join",
	}

	server.Start()
}
