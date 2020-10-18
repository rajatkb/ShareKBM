package windows

import (
	"fmt"
	"sharekbm/logger"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32DLL  = windows.NewLazyDLL("user32.dll")
	blockInput = user32DLL.NewProc("BlockInput")
)

//BlockInputFromHost ... blocks the input from the host
func BlockInputFromHost(logger *logger.Logger) {
	block := true
	r1, _, err := blockInput.Call(uintptr(unsafe.Pointer(&block)))
	if err != nil {
		logger.Warn("Input Stream from host is not blocked , cross clicks will be happening err:" + err.Error())
	} else {
		logger.Info(fmt.Sprintf("mouse status : %d", r1))
	}
}

//UnBlockInputFromHost ... unblocks the input
func UnBlockInputFromHost(logger *logger.Logger) {
	block := false
	r1, _, err := blockInput.Call(uintptr(unsafe.Pointer(&block)))
	if err != nil {
		logger.Warn("Input Stream from host is not blocked , cross clicks will be happening err:" + err.Error())
	} else {
		logger.Info(fmt.Sprintf("mouse status : %d", r1))
	}
}
