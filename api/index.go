package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kidboy-man/8-level-desent/app/config"
	v1 "github.com/kidboy-man/8-level-desent/app/controllers/http/v1"
	"github.com/kidboy-man/8-level-desent/app/repositories/inmemory"
	"github.com/kidboy-man/8-level-desent/app/services"
)

var (
	once   sync.Once
	engine *gin.Engine
)

func initEngine() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)
	engine = gin.Default()

	bookRepo := inmemory.NewBookRepository()
	bookService := services.NewBookService(bookRepo)
	authService := services.NewAuthService(cfg.JWTSecret)

	pingCtrl := v1.NewPingController()
	echoCtrl := v1.NewEchoController()
	authCtrl := v1.NewAuthController(authService)
	bookCtrl := v1.NewBookController(bookService)

	router := v1.NewRouter(engine, pingCtrl, echoCtrl, authCtrl, bookCtrl, authService)
	router.Setup()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(initEngine)
	engine.ServeHTTP(w, r)
}
