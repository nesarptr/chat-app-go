package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/nesarptr/chat-app-go/auth"
	"github.com/nesarptr/chat-app-go/controllers"
	"github.com/nesarptr/chat-app-go/socket"
)

func SetUpRoutes(app *fiber.App) {

	// unprotected routes

	// auth routes
	unProtected := app.Group("/")
	unProtected.Post("/signup", auth.SignUp)
	unProtected.Post("/signin", auth.SignIn)
	// ws routes
	ws := app.Group("/ws")
	ws.Use(socket.Upgrade)
	ws.Get("/", websocket.New(socket.WSHandler))

	// protected routes

	protected := app.Group("/auth", auth.Protected()...)
	protected.Get("/jwt", auth.Jwt)
	protected.Get("/users", controllers.GetUsers)
	protected.Get("/:from", controllers.GetMessages)

}
