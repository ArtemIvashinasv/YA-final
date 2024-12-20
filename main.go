package main

import (
	"log"

	"YandexPracticum-go-final-TODO/internal/server"
	"YandexPracticum-go-final-TODO/internal/server/handler"
	"YandexPracticum-go-final-TODO/internal/storage"

	"github.com/go-chi/chi"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Handle("/*", handler.GetFront())

	r.Post("/api/signin", handler.SignIn())

	r.Get("/api/nextdate", handler.GetNextDate)
	r.Post("/api/task", handler.Auth(handler.AddTask(db)))
	r.Get("/api/tasks", handler.Auth(handler.GetTasks(db)))
	r.Get("/api/task", handler.Auth(handler.GetTask(db)))
	r.Put("/api/task", handler.Auth(handler.UpdateTask(db)))
	r.Post("/api/task/done", handler.Auth(handler.DoneTask(db)))
	r.Delete("/api/task", handler.Auth(handler.DelTask(db)))

	server := new(server.Server)
	if err := server.Run(r); err != nil {
		log.Fatalf("Server can't start: %v", err)
		return
	}

	log.Println("Server stopped")
}
