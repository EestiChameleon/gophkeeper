package main

import (
	"github.com/EestiChameleon/gophkeeper/server/router"
	"github.com/EestiChameleon/gophkeeper/server/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// init the grpc server
	server, err := grpcserver.InitGRPCServer()
	if err != nil {
		log.Fatal(err)
	}

	// init storage
	if err = storage.Init(); err != nil {
		log.Fatal(err)
	}

	// channel to alert about shutdown
	gracefulShutdownChan := make(chan struct{})
	// channel to redirect the interrupt
	// we are looking after 3 syscall
	sigint := make(chan os.Signal, 3) // or size could be 1?
	// redirect registration
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// launch goroutine for received interrupt
	go func() {
		// we need only 1 signal to start the procedure
		<-sigint
		log.Println("server gracefully shutdown: start")
		if err = server.ShutDown(); err != nil {
			// ошибки закрытия Listener
			log.Printf("gRPC server shutdown err: %v", err)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(gracefulShutdownChan)
	}()

	// start the server
	if err = server.Start(); err != nil {
		log.Fatal(err)
	}

	// waiting the end of graceful shutdown procedure
	<-gracefulShutdownChan
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы

	if err = storage.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("server gracefully shutdown: done")
}
