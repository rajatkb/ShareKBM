package network

import (
	"fmt"
	"sharekbm/logger"

	"github.com/gorilla/websocket"
)

//Agent ... Represents an agent program or client
type Agent string

//Constants for agents type
const (
	CLIENTA Agent = "c"
	SERVERA Agent = "s"
)

//ServerManager ... client strcuture to manager incoming messages
type ServerManager struct {
	messages    map[int64](*chan []byte)
	logger      *logger.Logger
	connections []*websocket.Conn
}

//CreateServerManager ... creates server managaer
func CreateServerManager(logger *logger.Logger) *ServerManager {
	serverm := new(ServerManager)
	serverm.messages = make(map[int64]*chan []byte)
	serverm.logger = logger
	serverm.connections = make([]*websocket.Conn, 0)
	return serverm
}

//ConnectionHandler ... routine responsible for reading
func (server *ServerManager) ConnectionHandler(conn *websocket.Conn) {
	logger := server.logger

	defer conn.Close()
	conn.SetCloseHandler(func(code int, text string) error {
		logger.Warn(fmt.Sprintf("Closing with message %d , %s", code, text))
		return fmt.Errorf("Connection Closes with code : %d , text : %s ", code, text)
	})

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Something went wrong when reading from active socket connection")
		}
		if mt == websocket.TextMessage {
			fmt.Printf("server : %d %s", mt, string(p))
		}

		err2 := conn.WriteMessage(mt, p)

		if err2 != nil {
			conn.Close()
			logger.Fatal(fmt.Sprintf("Could not write message error : %s", err2.Error()))
		}
	}
}

//ClientManager manages client connection
type ClientManager struct {
	messages   chan []byte
	logger     *logger.Logger
	connection *websocket.Conn
}

//CreateClientManager ... creates client manager
func CreateClientManager(logger *logger.Logger) *ClientManager {
	clientm := new(ClientManager)
	clientm.messages = make(chan []byte)
	clientm.logger = logger
	return clientm
}

//ConnectionHandler ... manages client connections
func (client *ClientManager) ConnectionHandler(conn *websocket.Conn) error {
	logger := client.logger
	defer conn.Close()
	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("echo"))
		mt, p, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("error in reading from socker :%s", err.Error()))
			return err
		}
		if mt == websocket.TextMessage {
			fmt.Printf("client : %s", string(p))
		}
	}
}
