package main

import (
	"fmt"
	"os"
	logger "sharekbm/logger"
	network "sharekbm/network"
	"strconv"

	argparse "github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("sharekbm", "Controls operation parameter for sharekbm app")

	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "Port for application"})
	agentCode := parser.Selector("a", "agent", []string{string(network.SERVERA), string(network.CLIENTA)}, &argparse.Options{Required: true, Help: "Launch as clieant(c) or server(s)"})
	host := parser.String("t", "target", &argparse.Options{Required: true, Help: "Host address for Client(only used in Client)"})
	bufferSize := parser.Int("b", "buffer", &argparse.Options{Required: true, Help: "bufferSize"})

	clientName := parser.String("n", "name", &argparse.Options{Required: false, Default: "client-" + strconv.Itoa(os.Getpid()), Help: "an optional client name"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	logger := logger.GetLogger(*bufferSize)

	logger.Info(fmt.Sprintf(" Port used : %d for application", *port))
	logger.Info(fmt.Sprintf(" Agent : %s", *agentCode))
	logger.Info(fmt.Sprintf(" Host : %s", *host))
	logger.Info(fmt.Sprintf(" BufferSize : %d", *bufferSize))
	if *agentCode == string(network.CLIENTA) {
		logger.Info(fmt.Sprintf(" Client Name: %s", *clientName))
	}
	net := network.Network{
		Port:   port,
		Logger: logger,
		Host:   host,
	}

	if *agentCode == string(network.SERVERA) {
		agent := network.CreateServerManager(logger, bufferSize)
		net.CreateServer(agent)

	} else if *agentCode == string(network.CLIENTA) {
		agent := network.CreateClientManager(logger, bufferSize)
		net.CreateClient(agent)

	} else {
		logger.Fatal(fmt.Sprintf("Unknown agent m should be a ' %s ' or ' %s ' ", network.CLIENTA, network.SERVERA))
	}

	logger.Close()
}
