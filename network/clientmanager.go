package network

import (
	"fmt"
	"sharekbm/logger"

	"github.com/gorilla/websocket"
)

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
