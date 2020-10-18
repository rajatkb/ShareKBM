package device

import (
	"sharekbm/logger"
	network "sharekbm/network"

	hook "github.com/robotn/gohook"
)

//ControlManager ... creates the device server
type ControlManager struct {
	MouseInstance *MouseManager
}

func (controlServer *ControlManager) close() {
	controlServer.MouseInstance.close()
}

//CreateControlManager ... creates all device server , and starts the global event listener
func CreateControlManager(logger *logger.Logger, bufferSize int, agent network.Agent) *ControlManager {
	mouse := CreateMouseManager(logger, bufferSize, agent)

	// blocks for process to end
	go func() {
		s := hook.Start()
		<-hook.Process(s)
		// mm.blockOnListeners.Done()
	}()

	return &ControlManager{
		MouseInstance: mouse,
	}
}

//EventType ... described a device event type
type EventType int8

//MouseInteract , MouseMove , KeyPress ... contants for detecting event type
const (
	MouseInteract EventType = 0
	MouseMove     EventType = 1
	KeyPress      EventType = 2
)

//Event ... represents network type event
type Event interface {
	eventType() EventType
}
