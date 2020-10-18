package device

import (
	"sharekbm/logger"
	network "sharekbm/network"
	"sync"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

//MouseManager ... manages the input to the mouse using robotgo
type MouseManager struct {
	// mousePostition *MousePointer
	logger                 *logger.Logger
	agentCurrentPosition   chan *MousePointer
	agentChangePosition    chan *MousePointer
	agentListenClickEvent  chan *MouseClick
	agentRecieveClickEvent chan *MouseClick
	blockOnListeners       sync.WaitGroup
}

//CreateMouseManager ... creates a new MouseManeger
func CreateMouseManager(logger *logger.Logger, bufferSize int, agent network.Agent) *MouseManager {
	robotgo.SetMouseDelay(0) // setting 0 delay for recieving any mouse event
	mm := new(MouseManager)
	mm.logger = logger
	mm.agentCurrentPosition = make(chan *MousePointer, bufferSize)
	mm.agentChangePosition = make(chan *MousePointer, bufferSize)
	mm.agentListenClickEvent = make(chan *MouseClick, bufferSize)
	mm.agentRecieveClickEvent = make(chan *MouseClick, bufferSize)

	if agent == network.SERVERA {
		// mm.blockOnListeners.Add(1)
		mm.getCurrentPostion()
		mm.getCLickEvent()

	} else if agent == network.CLIENTA {
		go mm.setCurrentPosition()
		go mm.setClickEvent()
	}
	logger.Info("Mouse Manager Agent started")
	return mm
}

func (mouseManager *MouseManager) close() {
	close(mouseManager.agentCurrentPosition)
	close(mouseManager.agentChangePosition)
	close(mouseManager.agentListenClickEvent)
	close(mouseManager.agentRecieveClickEvent)
}

// GetCurrentPosReadStream ... get the stream of current position from host system
func (mouseManager *MouseManager) GetCurrentPosReadStream() <-chan *MousePointer {
	return mouseManager.agentCurrentPosition
}

// GetCurrentPosWriteStream ... input the location and the mouse position will be set accordingly
func (mouseManager *MouseManager) GetCurrentPosWriteStream() chan<- *MousePointer {
	return mouseManager.agentChangePosition
}

//GetClickListenStream ... get a stream of clicks happening in host system
func (mouseManager *MouseManager) GetClickListenStream() <-chan *MouseClick {
	return mouseManager.agentListenClickEvent
}

//GetClickRecieveStream ... get a stream of clicks happening in host system
func (mouseManager *MouseManager) GetClickRecieveStream() chan<- *MouseClick {
	return mouseManager.agentRecieveClickEvent
}

func (mouseManager *MouseManager) getCurrentPostion() {
	mouseManager.logger.Info("reading current mouse position for client")
	maxx, maxy := robotgo.GetScreenSize()
	var curx, cury int = robotgo.GetMousePos()

	hook.Register(hook.MouseMove, []string{}, func(ev hook.Event) {
		prevx, prevy := curx, cury
		curx, cury = int(ev.X), int(ev.Y)
		if (prevx-curx) == 0 && (prevy-cury) == 0 {
			return
		}
		mouseManager.agentCurrentPosition <- CreateMousePointer(int16(curx), int16(cury), int16(maxx), int16(maxy), int16(prevx), int16(prevy))
	})

}

func (mouseManager *MouseManager) setClickEvent() {

	mouseManager.logger.Info("ready for creating click events")
	for value := range mouseManager.agentRecieveClickEvent {
		switch value.MouseBtn {
		case LeftMouse:
			{
				if value.Press {
					robotgo.MouseToggle("down", "left")
				} else {
					robotgo.MouseToggle("up", "left")
				}
				continue
			}
		case RightMouse:
			{
				if value.Press {
					robotgo.MouseToggle("down", "right")
				} else {
					robotgo.MouseToggle("up", "right")
				}
				continue
			}
		case CenterMouse:
			{
				if value.Press {
					robotgo.MouseToggle("down", "center")
				} else {
					robotgo.MouseToggle("up", "center")
				}
				continue
			}
		case ScrollUpMouse:
			{
				robotgo.ScrollMouse(10, "up")
			}
		case ScrollDownMouse:
			{
				robotgo.ScrollMouse(10, "down")
			}
		}
	}
}

func (mouseManager *MouseManager) getCLickEvent() {
	mouseManager.logger.Info("reading current mouse event for client")

	var prev *MouseClick = nil

	doInsert := func(bt MouseButton, press bool) {
		if prev == nil || prev.MouseBtn != bt || prev.Press != press || bt == ScrollUpMouse || bt == ScrollDownMouse { // scroll up and scroll down must constantly need to be transmitted
			prev = CreateMouseClick(bt, press)
			mouseManager.agentListenClickEvent <- prev
		}
	}

	hook.Register(hook.MouseDown, []string{}, func(ev hook.Event) {
		if hook.MouseMap["left"] == ev.Button {
			doInsert(LeftMouse, true)
		} else if hook.MouseMap["right"] == ev.Button {
			doInsert(RightMouse, true)
		} else if hook.MouseMap["center"] == ev.Button {
			doInsert(CenterMouse, true)
		}
	})

	hook.Register(hook.MouseHold, []string{}, func(ev hook.Event) {
		if hook.MouseMap["left"] == ev.Button {
			doInsert(LeftMouse, true)
		} else if hook.MouseMap["right"] == ev.Button {
			doInsert(RightMouse, true)
		} else if hook.MouseMap["center"] == ev.Button {
			doInsert(CenterMouse, true)
		}
	})

	hook.Register(hook.MouseUp, []string{}, func(ev hook.Event) {
		if hook.MouseMap["left"] == ev.Button {
			doInsert(LeftMouse, false)
		} else if hook.MouseMap["right"] == ev.Button {
			doInsert(RightMouse, false)
		} else if hook.MouseMap["center"] == ev.Button {
			doInsert(CenterMouse, false)
		}
	})

	hook.Register(hook.MouseWheel, []string{}, func(ev hook.Event) {

		if ev.Rotation == 1 {
			doInsert(ScrollUpMouse, false)
		} else if ev.Rotation == -1 {
			doInsert(ScrollDownMouse, false)
		}
	})
}

func (mouseManager *MouseManager) setCurrentPosition() {
	mouseManager.logger.Info("Can set current mouse position for client")
	maxx, maxy := robotgo.GetScreenSize()
	var curx, cury int = robotgo.GetMousePos()
	for value := range mouseManager.agentChangePosition {
		value.adjustClickPosition(int16(maxx), int16(maxy)) // fixing the incoming data to thost system data
		if (curx-int(value.Curx)) == 0 && (cury-int(value.Cury)) == 0 {
			curx, cury = int(value.Curx), int(value.Cury)
			continue
		}
		curx, cury = int(value.Curx), int(value.Cury)
		robotgo.MoveMouse(int(value.Curx), int(value.Cury))
	}
}
