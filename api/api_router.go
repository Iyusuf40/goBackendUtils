package api

import (
	"github.com/Iyusuf40/goBackendUtils/api/controllers/user_controller"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServeAPI() {
	e := echo.New()
	e.Use(middleware.Recover())

	g := e.Group("/api")
	g.POST("/users", user_controller.SaveUser)
	g.GET("/users/:id", user_controller.GetUser)
	g.PUT("/users/:id", user_controller.UpdateUser)
	g.DELETE("/users/:id", user_controller.DeleteUser)

	// complete signup
	g.GET("/complete_signup/:signupId", user_controller.CompleteSignup)

	e.Logger.Fatal(e.Start(":" + config.ApiPort))
}
