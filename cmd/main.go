package main

import (
	"log"

	"github.com/LoaltyProgramm/to-do-app/internal/config"
	"github.com/LoaltyProgramm/to-do-app/internal/db"
	"github.com/LoaltyProgramm/to-do-app/internal/handlers"
	"github.com/LoaltyProgramm/to-do-app/internal/repository"
	"github.com/LoaltyProgramm/to-do-app/internal/server"
	"github.com/LoaltyProgramm/to-do-app/internal/service/authservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/cookieservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/middlewareservice"
	"github.com/LoaltyProgramm/to-do-app/internal/service/taskservice"
)

func main() {
	cfg, err := config.GetEnv()
	if err != nil {
		log.Println("Read env file complited")
	}
	
	if err := db.InitDB(cfg); err != nil {
		log.Fatal(err)
	}
	log.Println("Initialisation db file complited")

	DB := repository.NewTaskRepository(db.DB)
	taskService := service.NewTaskService(DB)
	authService := authservice.NewJWTService()
	cookieService := cookieservice.NewCookieService()
	middlewareService := middlewareservice.NewMiddlewareService(authService, cookieService, *cfg)
	taskHandlers := handlers.NewTaskHandlers(taskService, cfg, authService, cookieService, middlewareService)
	log.Println("Service start complited")

	taskHandlers.InitHandler()
	log.Println("Handlers start complited")

	log.Println("Server start complited")
	if err := server.StartServer(cfg); err != nil {
		log.Fatal(err)
	}
}	
