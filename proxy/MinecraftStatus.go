package proxy

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"log"
	"net"
	"strconv"
	"strings"
)

type MinecraftStatus struct {
	CurPlayers int
	MaxPlayers int
}

func GetServerStatus(serverIp string) (MinecraftStatus, error) {
	client, err := net.Dial("tcp", fmt.Sprintf("%s:25566", serverIp))
	if err != nil {
		return MinecraftStatus{}, err
	}

	client.Write([]byte{0xfe, 0x01})
	recvbuf := make([]byte, 128)
	read, err := client.Read(recvbuf)

	if err != nil {
		return MinecraftStatus{}, err
	}

	if read == 0 {
		return MinecraftStatus{}, errors.New("empty response")
	}

	decoder := unicode.UTF16(unicode.BigEndian, 0).NewDecoder()
	decodedData, _ := decoder.Bytes(recvbuf[0:read][3:])
	response := strings.Split(string(decodedData), "\x00")

	log.Println("Received response from the server")

	curPlayers, _ := strconv.ParseInt(response[4], 10, 32)
	maxPlayers, _ := strconv.ParseInt(response[5], 10, 32)
	return MinecraftStatus{
		CurPlayers: int(curPlayers),
		MaxPlayers: int(maxPlayers),
	}, nil
}
