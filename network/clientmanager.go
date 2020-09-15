package network

import (
	"fmt"
	"sharekbm/logger"
	"sync"

	"github.com/gorilla/websocket"
)

//ClientManager manages client connection
type ClientManager struct {
	incomingMessages chan []byte
	outgoingMessages chan []byte
	logger           *logger.Logger
	connection       *websocket.Conn
	bufferSize       *int
}

//CreateClientManager ... creates client manager
func CreateClientManager(logger *logger.Logger, bufferSize *int) *ClientManager {
	clientm := new(ClientManager)
	clientm.incomingMessages = make(chan []byte, *bufferSize)
	clientm.outgoingMessages = make(chan []byte, *bufferSize)
	clientm.logger = logger
	clientm.bufferSize = bufferSize
	return clientm
}

//GetSender ... give the channel for pushing in messages
func (clientm *ClientManager) GetSender() chan<- []byte {
	return clientm.outgoingMessages
}

//GetReceiver ... give the channel for reading messages
func (clientm *ClientManager) GetReceiver() <-chan []byte {
	return clientm.incomingMessages
}

func closeConnection(conn *websocket.Conn, logger *logger.Logger) {
	err := conn.WriteMessage(websocket.CloseAbnormalClosure, []byte("abonormal_closing"))
	if err != nil {
		logger.Error(fmt.Sprintf("count not send close message : %s", err.Error()))
	}
	err2 := conn.Close()
	if err2 != nil {
		logger.Error(fmt.Sprintf("Unable to close the socket err: %s", err2.Error()))
	}
}

//ConnectionHandler ... manages client connections
// returns error and sends a close message and closes the socker reosource
func (clientm *ClientManager) handleConnection(conn *websocket.Conn) error {
	logger := clientm.logger

	// conn.SetCloseHandler(func(c int, text string) error {
	// 	logger.Info(fmt.Sprintf("Server closed connection!! code : %d message: %s", c, text))
	// 	err := conn.Close()
	// 	if err != nil {
	// 		logger.Error(fmt.Sprintf("Connection to server already closed err: %s", err.Error()))
	// 	}
	// 	return fmt.Errorf("server closed connection")
	// })

	defer closeConnection(conn, logger)

	var wg sync.WaitGroup
	wg.Add(2)
	// writ send Messages and sending them to server
	go func() {
		for {
			select {
			case msg := <-clientm.outgoingMessages:
				{
					err := conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						logger.Error(fmt.Sprintf("could not write to host"))
						wg.Done()
						closeConnection(conn, logger)
						return
					}
				}
			default:
			}
		}
	}()

	// write read messages and sending to channel
	go func() {
		for {
			mt, p, err := conn.ReadMessage()
			if err != nil {
				logger.Error(fmt.Sprintf("error in reading from socker :%s", err.Error()))
				wg.Done()
				closeConnection(conn, logger)
				return
			}
			if mt == websocket.TextMessage {
				clientm.incomingMessages <- p
			}
		}
	}()

	wg.Wait()
	return fmt.Errorf("Socket connection died could not send or receive data")
}
