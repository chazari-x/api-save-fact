package main

import (
	"api-save-fact/config"
	"api-save-fact/model"
	"api-save-fact/server"
	"api-save-fact/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

func main() {
	// Чтение конфигурации
	cfg, err := config.ReadConfig("etc/config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	// Создание буфера
	buffer := make(chan model.Data, cfg.Buffer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewWorker(buffer, cfg)

	w.Start()

	s := server.NewServer(cfg, buffer)

	// Запуск сервера
	go func() {
		if err = s.Start(); err != nil {
			log.Println(err)
		}
	}()

	// Ожидание принудительного завершения
	SIG := make(chan os.Signal, 1)
	signal.Notify(SIG, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT)

	log.Println("Waiting for signal...")

	for {
		sig := <-SIG
		switch sig {
		case os.Kill, syscall.SIGKILL, syscall.SIGQUIT:
			log.Println("Received signal:", sig)
			return
		case syscall.SIGTERM, syscall.SIGINT:
			log.Println("Received signal:", sig)

			go func() {
				for {
					if sig := <-SIG; slices.Contains([]os.Signal{syscall.SIGKILL, syscall.SIGQUIT}, sig) {
						log.Println("Received signal:", sig)
						os.Exit(0)
					}
				}
			}()

			log.Println("Shutting down...")
			if err = s.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			w.Shutdown()

			return
		default:
			log.Println("Received signal:", sig)
		}
	}
}
