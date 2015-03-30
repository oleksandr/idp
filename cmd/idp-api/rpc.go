package main

import (
	"log"
	"os"

	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"github.com/oleksandr/idp/config"
	"github.com/oleksandr/idp/rpc"
	"github.com/oleksandr/idp/rpc/generated/services"
	"github.com/oleksandr/idp/usecases"
)

func startRPCServer(exitCh chan bool,
	domainInteractor usecases.DomainInteractor,
	userInteractor usecases.UserInteractor,
	sessionInteractor usecases.SessionInteractor) {

	addr := os.Getenv(config.EnvIDPRPCAddr)
	if addr == "" {
		addr = ":8001"
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTTransportFactory()
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Println("Error creating server transport: %s", err)
		exitCh <- true
		return
	}

	handler := rpc.NewAuthenticatorHandler()
	handler.DomainInteractor = domainInteractor
	handler.UserInteractor = userInteractor
	handler.SessionInteractor = sessionInteractor

	processor := services.NewAuthenticatorProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	log.Println("RPC API Server listening", addr)
	err = server.Serve()
	if err != nil {
		log.Println("Error in Serve: %s", err)
		exitCh <- true
	}
}
