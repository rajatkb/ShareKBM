package window

import (
	device "sharekbm/device"
	logger "sharekbm/logger"
	network "sharekbm/network"

	guuid "github.com/google/uuid"
)

//SystemWindow Handles the window data and related functionality
type SystemWindow struct {
	MaxHeight     int8
	MaxWidth      int8
	CurrentMouseX int8
	CurrentMouseY int8
}

//ServerWindoManager manages the bunch of windows coming in and going out
type ServerWindoManager struct {
	serverWindow           SystemWindow
	clientWindows          map[guuid.UUID]*SystemWindow
	serverIncomingMessages chan network.ServerMessage
	serverOutgoingMessages chan network.ServerMessage
	participants           **guuid.UUID
	clientWindowsGrid      []*SystemWindow
	maxClient              int
	clientCount            int
	deviceControlManager   *device.ControlManager
}

//CreateServerWindowManager .. creates a server manager window
func CreateServerWindowManager(maxClient int, logger *logger.Logger, bufferSize int, deviceControlManager *device.ControlManager) *ServerWindoManager {
	serverWindoManager := new(ServerWindoManager)
	serverWindoManager.clientWindows = make(map[guuid.UUID]*SystemWindow)
	serverWindoManager.serverIncomingMessages = make(chan network.ServerMessage)
	serverWindoManager.serverOutgoingMessages = make(chan network.ServerMessage)
	serverWindoManager.deviceControlManager = deviceControlManager
	serverWindoManager.clientWindowsGrid = make([]*SystemWindow, maxClient)
	serverWindoManager.maxClient = maxClient
	serverWindoManager.clientCount = 0

	return serverWindoManager
}

// go routine responsible for reading device input streams and pushing it to correct client
// all server to client switches takes place at this point
func (serverWindowManager *ServerWindoManager) processDeviceMessages() {

}

//GetWindowServerWriteStream ... give u a writable stream of incoming message
func (serverWindowManager *ServerWindoManager) GetWindowServerWriteStream() chan<- network.ServerMessage {
	return serverWindowManager.serverIncomingMessages
}

//GetWindowServerReadStream ... give you a readable stream of outgoing messages from the server
func (serverWindowManager *ServerWindoManager) GetWindowServerReadStream() <-chan network.ServerMessage {
	return serverWindowManager.serverOutgoingMessages
}

const connectEvent device.EventType = 100

// messages coming from socket server into the window server
func (serverWindowManager *ServerWindoManager) processIncomingMessage() {
	// for value := range serverWindowManager.serverIncomingMessages {

	// }
}
