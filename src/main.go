package main

import (
	"fmt"
	"twometer.dev/craftignite/minecraft"
)

func main() {
	fmt.Println("CraftIgnite starting up")
	server := minecraft.Server{}
	server.Start()
}
