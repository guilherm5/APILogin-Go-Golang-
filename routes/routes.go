package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guilherm5/GoAndReact/controllers"
)

func Setup(app *fiber.App) error {

	//registrar usuario
	app.Post("/api/register", controllers.Register)
	//logar usuario e pegar jwt
	app.Post("/api/login", controllers.Login)
	//ver qual usuario esta logado
	app.Get("/api/user", controllers.User)
	//deslogando usuario
	app.Post("/api/logout", controllers.Logout)

	return nil
}
