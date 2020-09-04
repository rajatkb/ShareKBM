package network

import (
	"fmt"
	"net/http"
	"net/url"
	loglib "sharekbm/logger"

	"github.com/gorilla/websocket"
)

type websocketStatus int8

// SUCCESS , FAIL STATUS
const (
	SUCCESS websocketStatus = 1
	FAILED  websocketStatus = -1
)

func dummyReader(conn *websocket.Conn, logger *loglib.Logger) websocketStatus {

	defer conn.Close()

	conn.SetCloseHandler(func(code int, text string) error {
		logger.Warn(fmt.Sprintf("Closing with message %d , %s", code, text))
		return fmt.Errorf("Connection Closes with code : %d , text : %s ", code, text)
	})

	for {
		mt, p, err := conn.ReadMessage()
		if err != nil {
			logger.Error("Something went wrong when reading from active socket connection")
			return FAILED
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

func wsEndPoint(logger *loglib.Logger, serverm *ServerManager) func(w http.ResponseWriter, r *http.Request) {

	var websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocketUpgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to start web socket error: %s", err.Error()))
		}

		serverm.handleConnection(conn) // non blocking call to handle a connection
		return
	}

}

func notfound(logger *loglib.Logger) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("unknown url hit")
		fmt.Fprintf(w, "No page here pal ¯\\_(ツ)_/¯")
	}
}

func config(logger *loglib.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		go logger.Info("Config page hit")
		fmt.Fprintf(w, "config page")
	}
}

//Network ... Structure holding all network related functionality
type Network struct {
	Port   *int
	Logger *loglib.Logger
	Host   *string
}

// CreateServer ... provides a wss server
func (network *Network) CreateServer(serverm *ServerManager) {
	logger := network.Logger

	http.HandleFunc("/config", config(network.Logger))
	http.HandleFunc("/wss", wsEndPoint(network.Logger, serverm))
	http.HandleFunc("/", notfound(network.Logger))

	logger.Info(fmt.Sprintf("Server starting in port :%d", *network.Port))

	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", *network.Port), nil)

	if err != nil {
		logger.Fatal(err.Error())
	}
}

// CreateClient ... Provides a wss client for sending message to the server
func (network *Network) CreateClient(client *ClientManager) {
	logger := network.Logger
	logger.Info(fmt.Sprintf("Client starting in port :%d", *network.Port))

	for {
		logger.Info("Client attempting to connect !!")
		addr := fmt.Sprintf("%s:%d", *network.Host, *network.Port)
		url := url.URL{Scheme: "ws", Host: addr, Path: "/wss"}
		conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
		logger.Info("Client established a connection !!")
		if err != nil {
			logger.Error(fmt.Sprintf("Failed when creating connection error : %s ", err.Error()))
			return
		}
		// reader domain
		err2 := client.ConnectionHandler(conn) // blocking call to handle a connection
		if err2 != nil {
			logger.Error("Disconnected from Server, retrying connection")
		}
	}

}
