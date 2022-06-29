package main

import (
	"net"

	"github.com/badfan/inno-taxi-order-service/app"
	"github.com/badfan/inno-taxi-order-service/app/resources"
	"github.com/badfan/inno-taxi-order-service/app/rpc"
	"github.com/badfan/inno-taxi-order-service/app/services/order"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	sugarLogger := logger.Sugar()
	return sugarLogger
}

func InitGRPCClients(rpcConfig *app.RPCConfig, logger *zap.SugaredLogger) (*grpc.ClientConn, *grpc.ClientConn) {
	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	userClientConn, err := grpc.Dial("localhost:"+rpcConfig.RPCUserPort, options...)
	if err != nil {
		logger.Fatalf("error occurred while connecting to user GRPC server: %s", err.Error())
	}

	driverClientConn, err := grpc.Dial("localhost:"+rpcConfig.RPCDriverPort, options...)
	if err != nil {
		logger.Fatalf("error occurred while connecting to driver GRPC server: %s", err.Error())
	}

	return userClientConn, driverClientConn
}

func InitGRPCServer(rpcService *order.RPCService, rpcConfig *app.RPCConfig, logger *zap.SugaredLogger) {
	listener, err := net.Listen("tcp", "localhost:"+rpcConfig.RPCOrderPort)
	if err != nil {
		logger.Fatalf("failed to up grpc server: %s", err.Error())
	}

	var options []grpc.ServerOption
	rpcServer := grpc.NewServer(options...)
	rpc.RegisterOrderServiceServer(rpcServer, rpcService.OrderServiceServer)
	rpcServer.Serve(listener)
}

func main() {
	logger := InitLogger()
	defer logger.Sync()

	rpcConfig, err := app.NewRPCConfig()
	if err != nil {
		logger.Fatalf("rpcconfig error: %s", err.Error())
	}
	dbConfig, err := app.NewDBConfig()
	if err != nil {
		logger.Fatalf("dbconfig error: %s", err.Error())
	}

	resource, err := resources.NewResource(dbConfig, logger)
	if err != nil {
		logger.Fatalf("db error: %s", err.Error())
	}
	defer resource.Db.Close()

	userClientConn, driverClientConn := InitGRPCClients(rpcConfig, logger)
	defer userClientConn.Close()
	defer driverClientConn.Close()

	rpcService := order.NewRPCService(userClientConn, driverClientConn, resource, rpcConfig, logger)

	InitGRPCServer(rpcService, rpcConfig, logger)
}
