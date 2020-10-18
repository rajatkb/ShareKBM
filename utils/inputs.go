package utils

import (
	"runtime"
	"sharekbm/logger"
	windows "sharekbm/utils/windows"
)

//BlockInputFromHost ... blocks the input from the host
func BlockInputFromHost(logger *logger.Logger) bool {
	if runtime.GOOS == "windows" {
		windows.BlockInputFromHost(logger)
		return true
	}
	if runtime.GOOS == "unix" {
		logger.Warn("Unsupported method , cross clicks happening")
		return false
	}
	return false
}

//UnBlockInputFromHost ... unblocks the input
func UnBlockInputFromHost(logger *logger.Logger) bool {
	if runtime.GOOS == "windows" {
		windows.UnBlockInputFromHost(logger)
		return true
	}
	if runtime.GOOS == "unix" {
		logger.Warn("Unsupported method , cross clicks happening")
		return false
	}
	return false
}
