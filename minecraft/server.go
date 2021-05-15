package minecraft

import (
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"net"
)

type Server struct {
	listener        net.Listener
	Motd            string
	KickMessage     string
	TooltipMessage  string
	VersionName     string
	MaxPlayerCount  int
	ConnectCallback func()
}

type Client struct {
	conn            net.Conn
	currentState    int
	protocolVersion int
}

const stateJson = `
{
    "version": {
        "name": "%s",
        "protocol": %d
    },
    "players": {
        "max": %d,
        "online": 0,
        "sample": [
            {
                "name": "%s",
                "id": "00000000-0000-0000-0000-000000000000"
            }
        ]
    },
    "description": {
        "text": "%s"
    }
}`

func (server *Server) Start() {
	var err error
	server.listener, err = net.Listen("tcp", "localhost:25565")

	if err != nil {
		fmt.Println("Server failed to start: ", err.Error())
		return
	}

	for {
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Server failed to accept client: " + err.Error())
			return
		}

		go server.handleClient(conn)
	}
}

func (server *Server) Stop() {
	err := server.listener.Close()
	if err != nil {
		fmt.Println("Server failed to stop: " + err.Error())
		return
	}
}

func (server *Server) handleClient(conn net.Conn) {
	receiveBuf := make([]byte, 1024)
	client := Client{conn, 0, 0}

	for {
		read, err := conn.Read(receiveBuf)
		if err != nil {
			return
		}

		packet := Buffer{receiveBuf[0:read], 0}
		if packet.data[0] == 0xfe {
			server.handleLegacyPing(client, &packet)
			continue
		}

		packet.ReadVarInt()
		pid := packet.ReadVarInt()

		switch {
		case client.currentState == 0 && pid == 0:
			client.protocolVersion = packet.ReadVarInt()
			packet.Skip(packet.ReadVarInt() + 2)
			client.currentState = packet.ReadVarInt()
		case client.currentState == 1:
			server.handleStatusPacket(client, pid, &packet)
		case client.currentState == 2:
			server.handleLoginPacket(client, pid, &packet)
		}
	}
}

func (server *Server) handleLegacyPing(client Client, packet *Buffer) {
	response := Buffer{make([]byte, 1024), 0}
	encoder := unicode.UTF16(unicode.BigEndian, 0).NewEncoder()
	infoString := fmt.Sprintf("§1\x00127\x00%s\x00%s\x000\x00%d", server.VersionName, server.Motd, server.MaxPlayerCount)
	utf16be, _ := encoder.String(infoString)

	response.WriteByte(0xFF)
	response.WriteShortBE(uint16(len(infoString) - 1))
	response.WriteBytes([]byte(utf16be))
	sendPacket(client.conn, &response)
}

func (server *Server) handleStatusPacket(client Client, pid int, packet *Buffer) {
	switch pid {
	case 0: // Status Request
		response := Buffer{make([]byte, 1024), 0}
		response.WriteVarInt(0) // Status Response
		response.WriteString(fmt.Sprintf(stateJson, server.VersionName, client.protocolVersion, server.MaxPlayerCount, server.TooltipMessage, server.Motd))
		sendPacket(client.conn, &response)
	case 1: // Ping
		response := Buffer{make([]byte, 16), 0}
		response.WriteVarInt(1) // Pong
		response.WriteLong(packet.ReadLong())
		sendPacket(client.conn, &response)
	}
}

func (server *Server) handleLoginPacket(client Client, pid int, packet *Buffer) {
	switch pid {
	case 0: // Login Request
		response := Buffer{make([]byte, 128), 0}
		response.WriteVarInt(0) // Login Disconnect
		response.WriteString(fmt.Sprintf(`{ "text": "%s" }`, server.KickMessage))
		sendPacket(client.conn, &response)
		server.ConnectCallback()
	}
}

func sendPacket(conn net.Conn, packet *Buffer) {
	container := Buffer{make([]byte, len(packet.data)+8), 0}
	container.WriteVarInt(int(packet.offset))
	container.WriteBytes(packet.data[0:packet.offset])
	_, err := conn.Write(container.data[0:container.offset])
	if err != nil {
		fmt.Println("Failed to send packet to client")
	}
}
