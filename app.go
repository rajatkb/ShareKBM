package main

import (
	"fmt"
	"os"
	logger "sharekbm/logger"
	network "sharekbm/network"

	argparse "github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("sharekbm", "Controls operation parameter for sharekbm app")

	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "Port for application"})
	agent := parser.Selector("a", "agent", []string{string(network.SERVERA), string(network.CLIENTA)}, &argparse.Options{Required: true, Help: "Launch as clieant or server"})
	host := parser.String("t", "target", &argparse.Options{Required: true, Help: "Host address for Client(only used in Client)"})
	bufferSize := parser.Int("b", "buffer", &argparse.Options{Required: true, Help: "bufferSize"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	logger := logger.GetLogger(*bufferSize)

	logger.Info(fmt.Sprintf(" Port used : %d for application", *port))
	logger.Info(fmt.Sprintf(" Agent : %s", *agent))
	logger.Info(fmt.Sprintf(" Host : %s", *host))
	logger.Info(fmt.Sprintf(" BufferSize : %d", *bufferSize))

	net := network.Network{
		Port:   port,
		Logger: logger,
		Host:   host,
	}

	if *agent == string(network.SERVERA) {
		agnt := network.CreateServerManager(logger)
		net.CreateServer(agnt)

	} else if *agent == string(network.CLIENTA) {
		agnt := network.CreateClientManager(logger)
		net.CreateClient(agnt)

	} else {
		logger.Fatal(fmt.Sprintf("Unknown agent m should be a ' %s ' or ' %s ' ", network.CLIENTA, network.SERVERA))
	}

	logger.Close()
}
