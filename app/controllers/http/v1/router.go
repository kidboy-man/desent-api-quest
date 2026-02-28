package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kidboy-man/8-level-desent/app/controllers/middlewares"
	"github.com/kidboy-man/8-level-desent/app/services"
)

type Router struct {
	engine         *gin.Engine
	pingController *PingController
	echoController *EchoController
	authController *AuthController
	bookController *BookController
	authService    *services.AuthService
}

func NewRouter(
	engine *gin.Engine,
	pingCtrl *PingController,
	echoCtrl *EchoController,
	authCtrl *AuthController,
	bookCtrl *BookController,
	authService *services.AuthService,
) *Router {
	return &Router{
		engine:         engine,
		pingController: pingCtrl,
		echoController: echoCtrl,
		authController: authCtrl,
		bookController: bookCtrl,
		authService:    authService,
	}
}

func (r *Router) Setup() {
	r.engine.GET("/ping", r.pingController.Ping)
	r.engine.POST("/echo", r.echoController.Echo)
	r.engine.POST("/auth/token", r.authController.GenerateToken)

	protected := r.engine.Group("/")
	protected.Use(middlewares.AuthMiddleware(r.authService))
	{
		protected.POST("/books", r.bookController.Create)
		protected.GET("/books", r.bookController.GetAll)
		protected.GET("/books/:id", r.bookController.GetByID)
		protected.PUT("/books/:id", r.bookController.Update)
		protected.DELETE("/books/:id", r.bookController.Delete)
	}
}
