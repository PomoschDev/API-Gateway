package main

import (
	_ "apiGateway/docs"
	"apiGateway/iternal/grpc"
	"apiGateway/iternal/server"
	"apiGateway/pkg/config"
	"apiGateway/pkg/logger"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	if err := logger.New(); err != nil {
		panic(any(fmt.Errorf("Ошибка при инициализации логера: %v\n", err)))
	}

	//Инициализация конфигурации
	cfg := config.MustLoad()

	//Соединение с сервисом DatabaseService
	grpc_client := grpc.New(cfg)

	//Создание сервера маршрутизации
	srv := server.New(cfg, grpc_client)

	//Правильное завершение сервиса
	{
		wait := time.Second * 15

		// Запуск сервера в отдельном потоке
		go func() {
			logger.Info("Сервер запущен на адресе: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				logger.Error("Ошибка при прослушивании сервера: %v", err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c

		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		grpc_client.Connect.Close()
		srv.Shutdown(ctx)
		srv.Close()
		logger.Warn("Выключение сервера")
		os.Exit(0)
	}
}
