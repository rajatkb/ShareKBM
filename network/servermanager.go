package network

import (
	"fmt"
	"sharekbm/logger"

	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//Client ... client strcuture to be managed by Server for writing messages into it
type Client struct {
	writer chan []byte
	reader chan []byte
	conn   *websocket.Conn
	id     guuid.UUID
}

func (client *Client) readMessageStream(serverManager *ServerManager) {
	for {
		_, p, err := client.conn.ReadMessage()
		if err != nil {
			serverManager.logger.Error(fmt.Sprintf("Error occureued wher reading from client id: %s , closing connection for same", client.id.String()))
			client.conn.Close()
			return
		}
		client.reader <- p
		serverManager.
	}
}

//ServerManager ... client strcuture to manager incoming messages
type ServerManager struct {
	clients      map[guuid.UUID]*Client
	incomingPipe chan struct {
		uuid    guuid.UUID
		message []byte
	}
	outgoingPipe chan struct {
		uuid    guuid.UUID
		message []byte
	}
	onConnection func(client *Client)
	logger       *logger.Logger
}

func (serverManager *ServerManager) send2IncomingPipe(id *guuid.UUID, payload *[]byte) {
	data := struct {
		uuid    guuid.UUID
		message byte
	}{
		uuid: *id, 
		message: *payload
	}
	serverManager.incomingPipe <- data
}

//CreateServerManager ... creates server managaer
func CreateServerManager(logger *logger.Logger) *ServerManager {
	serverm := new(ServerManager)
	serverm.clients = make(map[guuid.UUID]*Client)
	serverm.incomingPipe = make(chan struct {
		uuid    guuid.UUID
		message []byte
	})
	serverm.outgoingPipe = make(chan struct {
		uuid    guuid.UUID
		message []byte
	})
	serverm.logger = logger
	return serverm
}

//createClient creates client data structure to manage individual clients in the server
func (server *ServerManager) createClient(conn *websocket.Conn) {
	logger := server.logger
	client := new(Client)
	// assigning memory location with defaults
	*client = Client{writer: make(chan []byte), reader: make(chan []byte), conn: conn, id: guuid.New()}
	server.clients[client.id] = client
	logger.Info(fmt.Sprintf("New Client Joined with id : %s", client.id.String()))

}

//handleConnection ... routine responsible for reading
func (server *ServerManager) handleConnection(conn *websocket.Conn) {
	logger := server.logger

	conn.SetCloseHandler(func(code int, text string) error {
		logger.Warn(fmt.Sprintf("Closing with message : %d , %s", code, text))
		return fmt.Errorf("Connection Closes with code : %d , text : %s ", code, text)
	})

	// for {
	// 	mt, p, err := conn.ReadMessage()
	// 	if err != nil {
	// 		logger.Error(fmt.Sprintf("Coult not read to server : %s", err.Error()))
	// 		break
	// 	}
	// 	if mt == websocket.TextMessage {
	// 		fmt.Printf("server : %d %s", mt, string(p))
	// 	}

	// 	err2 := conn.WriteMessage(mt, p)

	// 	if err2 != nil {
	// 		logger.Error(fmt.Sprintf("Could not write message error : %s", err2.Error()))
	// 		break
	// 	}
	// }
}
