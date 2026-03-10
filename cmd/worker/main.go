package main

import (
	"log"

	"github.com/shutsuensha/go-tasks/internal/config"
	"github.com/shutsuensha/go-tasks/internal/queue"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	go queue.StartScheduler(cfg.RedisAddr)

	worker := queue.NewWorker(cfg.RedisAddr)
	handlers := queue.NewHandlers()

	log.Println("worker started")

	if err := worker.Run(handlers); err != nil {
		log.Fatal(err)
	}
}