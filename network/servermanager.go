package network

import (
	"fmt"
	"sharekbm/logger"
	"sync"

	guuid "github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//Client ... client strcuture to be managed by Server for writing messages into it
type Client struct {
	// inputStream  chan []byte
	outputStream chan []byte
	conn         *websocket.Conn
	id           guuid.UUID
}

//GetID ... gives the uuid of a connection
func (client *Client) GetID() *guuid.UUID {
	return &client.id
}

//GetConnection ... gives connection for each client
func (client *Client) GetConnection() *websocket.Conn {
	return client.conn
}

// //GetInputStream ... gives a readable stream from where client information can be read in
// func (client *Client) GetInputStream() <-chan []byte {
// 	return client.inputStream
// }

//GetOutputStream ... gives a writable stream where client specific information can be sent
func (client *Client) GetOutputStream() chan<- []byte {
	return client.outputStream
}

//createClient creates client data structure to manage individual clients in the server
func (server *ServerManager) createClient(conn *websocket.Conn) *Client {
	logger := server.logger
	client := new(Client)
	// assigning memory location with defaults
	// assigning by memory copy
	*client = Client{
		//  inputStream: make(chan []byte),
		outputStream: make(chan []byte, *server.bufferSize),
		conn:         conn,
		id:           guuid.New()}

	server.clientsLock.Lock()
	server.clients[client.id] = client
	server.clientsLock.Unlock()

	logger.Info(fmt.Sprintf("New Client Joined with id : %s", client.id.String()))

	return client
}

func (server *ServerManager) destroyClient(id *guuid.UUID) {
	logger := server.logger
	client, ok := server.clients[*id]
	if ok {
		// close(client.inputStream)
		close(client.outputStream)
		err := client.conn.WriteMessage(websocket.CloseMessage, []byte("closing"))
		if err != nil {
			logger.Error(fmt.Sprintf("count not send closing message to clint : %s", *id))
		}
		err2 := client.conn.Close()
		if err2 != nil {
			logger.Error(fmt.Sprintf("count not close properly connection to : %s", *id))
		}

		logger.Info(fmt.Sprintf("Client with id : %s disconnected", *id))
		server.clientsLock.Lock()
		delete(server.clients, *id)
		server.clientsLock.Unlock()
		return
	}
	logger.Info(fmt.Sprintf("Something is wrong client of id: %s , not found ", id.String()))

}

//ServerMessage ... structure for encloding messages in channel
type ServerMessage struct {
	uuid    guuid.UUID
	message []byte
}

//GetID ... returns the unique id of the servermessage
func (sm *ServerMessage) GetID() *guuid.UUID {
	return &sm.uuid
}

//GetMessage ... returns the binary message of the servermessage
func (sm *ServerMessage) GetMessage() *[]byte {
	return &sm.message
}

//ServerManager ... client strcuture to manager incoming messages
type ServerManager struct {
	clientsLock  sync.Mutex
	clients      map[guuid.UUID]*Client
	incomingPipe chan ServerMessage
	outgoingPipe chan ServerMessage
	// onConnection func(client *Client)
	logger     *logger.Logger
	bufferSize *int
}

//GetOutStream ... give a writable stream to sending messaged to clients
func (server *ServerManager) GetOutStream() chan<- ServerMessage {
	return server.outgoingPipe
}

//GetInStream ... give a readable stream for reading messages from clients
func (server *ServerManager) GetInStream() <-chan ServerMessage {
	return server.incomingPipe
}

//CreateServerManager ... creates server managaer
func CreateServerManager(logger *logger.Logger, bufferSize *int) *ServerManager {
	serverm := new(ServerManager)
	serverm.bufferSize = bufferSize
	serverm.clients = make(map[guuid.UUID]*Client)
	serverm.incomingPipe = make(chan ServerMessage, *bufferSize)
	serverm.outgoingPipe = make(chan ServerMessage, *bufferSize)
	serverm.logger = logger

	go serverm.readOutgoingPipe()

	return serverm
}

// read input from server and pass it to the selected client
func (server *ServerManager) readOutgoingPipe() {
	// blocking loop until server is closed down or destroyed
	for value := range server.outgoingPipe {
		msg := value.message
		id := value.uuid
		go func() {
			server.clientsLock.Lock()
			client, ok := server.clients[id]
			server.clientsLock.Unlock()
			if ok {
				client.GetOutputStream() <- msg
			}
		}()
	}
}

//handleConnection ... routine responsible for reading
func (server *ServerManager) handleConnection(conn *websocket.Conn) {
	logger := server.logger
	// conn.SetCloseHandler(func(code int, text string) error {
	// 	logger.Warn(fmt.Sprintf("Closing with message : %d , %s", code, text))
	// 	return fmt.Errorf("Connection Closes with code : %d , text : %s ", code, text)
	// })

	client := server.createClient(conn)

	var clientWait sync.WaitGroup
	clientWait.Add(2)
	//read from outputStream of client and push it to socket
	go func() {

		for value := range client.outputStream {
			err := conn.WriteMessage(websocket.TextMessage, value)
			if err != nil {
				logger.Error(fmt.Sprintf("could not write to host err :%s", err.Error()))
				server.destroyClient(&client.id)
				clientWait.Done()
			}
		}
	}()
	// read from client stream and pass it into server
	go func() {
		for {
			mt, p, err := client.conn.ReadMessage()
			if err != nil {
				logger.Error(fmt.Sprintf("error in reading from socket :%s for id : %s ", err.Error(), client.GetID().String()))
				server.destroyClient(&client.id)
				clientWait.Done()
				return
			}
			if mt == websocket.TextMessage {
				server.incomingPipe <- ServerMessage{uuid: *client.GetID(), message: p}
			} else if mt == websocket.CloseMessage {
				logger.Info(fmt.Sprintf("Recieved closing message from client id: %s , closing from server end", client.id.String()))
				server.destroyClient(&client.id)
				return
			}
		}
	}()

	clientWait.Wait()
	server.destroyClient(&client.id)
}
