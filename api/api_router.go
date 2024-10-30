package api

import (
	"net/http"

	"github.com/Iyusuf40/goBackendUtils/api/controllers/user_controller"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e = echo.New()
var g = e.Group("/api")

func GetApiGroup() *echo.Group {
	return g
}

func ServeAPI() {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	g.POST("/users", user_controller.SaveUser)
	g.GET("/users/:id", user_controller.GetUser)
	g.PUT("/users/:id", user_controller.UpdateUser)
	g.DELETE("/users/:id", user_controller.DeleteUser)

	// complete signup
	g.GET("/complete_signup/:signupId", user_controller.CompleteSignup)

	e.Logger.Fatal(e.Start(":" + config.ApiPort))
}
