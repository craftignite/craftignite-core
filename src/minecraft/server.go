package minecraft

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
}

type Client struct {
	conn            net.Conn
	currentState    int
	protocolVersion int
}

const kickMessage = "§c[CraftIgnite]§r The server is currently starting.\nPlease try to reconnect in a minute."
const stateJson = `
{
    "version": {
        "name": "9.8.7",
        "protocol": %d
    },
    "players": {
        "max": 100,
        "online": 101,
        "sample": [
            {
                "name": "test",
                "id": "00000000-0000-0000-0000-000000000000"
            }
        ]
    },
    "description": {
        "text": "§eCraftIgnite Minecraft Proxy\n§8Server currently sleeping"
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

		go handleClient(conn)
	}
}

func (server *Server) Stop() {
	err := server.listener.Close()
	if err != nil {
		fmt.Println("Server failed to stop: " + err.Error())
		return
	}
}

func handleClient(conn net.Conn) {
	receiveBuf := make([]byte, 1024)
	client := Client{conn, 0, 0}

	for {
		read, err := conn.Read(receiveBuf)
		if err != nil {
			return
		}

		packet := Buffer{receiveBuf[0:read], 0}
		if packet.data[0] == 0xfe {
			handleLegacyPing(client, &packet)
			continue
		}

		len := packet.ReadVarInt()
		pid := packet.ReadVarInt()
		fmt.Printf("Received packet #%d (%d bytes)\n", pid, len)

		switch {
		case client.currentState == 0 && pid == 0:
			client.protocolVersion = packet.ReadVarInt()
			packet.Skip(packet.ReadVarInt() + 2)
			client.currentState = packet.ReadVarInt()
		case client.currentState == 1:
			handleStatusPacket(client, pid, &packet)
		case client.currentState == 2:
			handleLoginPacket(client, pid, &packet)
		}
	}
}

func handleLegacyPing(client Client, packet *Buffer) {
	fmt.Println("Legacy ping received")
}

func handleStatusPacket(client Client, pid int, packet *Buffer) {
	switch pid {
	case 0: // Status Request
		response := Buffer{make([]byte, 1024), 0}
		response.WriteVarInt(0) // Status Response
		response.WriteString(fmt.Sprintf(stateJson, client.protocolVersion))
		sendPacket(client.conn, &response)
	case 1: // Ping
		response := Buffer{make([]byte, 16), 0}
		response.WriteVarInt(1) // Pong
		response.WriteLong(packet.ReadLong())
		sendPacket(client.conn, &response)
	}
}

func handleLoginPacket(client Client, pid int, packet *Buffer) {
	switch pid {
	case 0: // Login Request
		response := Buffer{make([]byte, 128), 0}
		response.WriteVarInt(0) // Login Disconnect
		response.WriteString(fmt.Sprintf(`{ "text": "%s" }`, kickMessage))
		sendPacket(client.conn, &response)
	}
}

func sendPacket(conn net.Conn, packet *Buffer) {
	container := Buffer{make([]byte, len(packet.data)+16), 0}
	container.WriteVarInt(int(packet.offset))
	container.WriteBytes(packet.data[0:packet.offset])
	_, err := conn.Write(container.data[0:container.offset])
	if err != nil {
		fmt.Println("Failed to send packet to client")
	}
}
