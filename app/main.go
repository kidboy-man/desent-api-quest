package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kidboy-man/8-level-desent/app/config"
	v1 "github.com/kidboy-man/8-level-desent/app/controllers/http/v1"
	"github.com/kidboy-man/8-level-desent/app/repositories/inmemory"
	"github.com/kidboy-man/8-level-desent/app/services"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)
	engine := gin.Default()

	bookRepo := inmemory.NewBookRepository()

	bookService := services.NewBookService(bookRepo)
	authService := services.NewAuthService(cfg.JWTSecret)

	pingCtrl := v1.NewPingController()
	echoCtrl := v1.NewEchoController()
	authCtrl := v1.NewAuthController(authService)
	bookCtrl := v1.NewBookController(bookService)

	router := v1.NewRouter(engine, pingCtrl, echoCtrl, authCtrl, bookCtrl, authService)
	router.Setup()

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
