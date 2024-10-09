package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Hello world!")

	//Правильное завершение сервиса
	{
		wait := time.Second * 15

		// Запуск сервера в отдельном потоке

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		<-c

		_, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		os.Exit(0)
	}
}
