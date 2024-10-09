package grpc

import (
	DatabaseServicev1 "apiGateway/iternal/DatabaseService"
	"apiGateway/pkg/config"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
)

type Api struct {
	Cfg     *config.Config                          // Конфигурация приложения
	Client  DatabaseServicev1.DatabaseServiceClient // Клиент для взаимодействия с gRPC-сервером
	Connect *grpc.ClientConn                        //Соединение с grpc
}

func New(cfg *config.Config) *Api {
	cc, err := grpc.NewClient(grpcAddress(cfg), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(any(fmt.Errorf("grpc server connection failed: %v", err)))
	}

	return &Api{
		Cfg:     cfg,
		Client:  DatabaseServicev1.NewDatabaseServiceClient(cc),
		Connect: cc,
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.GRPCServer.Host, strconv.Itoa(cfg.GRPCServer.Port))
}
