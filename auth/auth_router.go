package auth

import (
	"fmt"
	"net/http"

	"github.com/Iyusuf40/goBackendUtils/api/controllers"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServeAUTH() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	g := e.Group("/auth")
	g.POST("/login", Login)
	g.GET("/logout", Logout)

	g.PUT("/isloggedin", IsLoggedIn)

	g.POST("forgot_password", ForgotPassword)
	g.POST("reset_password/:passwordResetToken", ResetPassword)

	e.Logger.Fatal(e.Start(":" + config.AuthPort))
}

var AUTH_HANDLER = MakeAuthHandler(config.TempStoreDb,
	config.UsersDatabase, config.UsersRecords)

func Login(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	email, email_ok := userDesc["email"].(string)

	if !email_ok {
		response["error"] = "email required to login"
		return c.JSON(http.StatusBadRequest, response)
	}

	password, password_ok := userDesc["password"].(string)

	if !password_ok {
		response["error"] = "password required to login"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId := AUTH_HANDLER.HandleLogin(email, password)

	if sessionId == "" {
		response["error"] = "failed to login"
		return c.JSON(http.StatusBadRequest, response)
	}
	response["sessionId"] = sessionId
	return c.JSON(http.StatusOK, response)
}

func Logout(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	sessionDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId, sessionId_ok := sessionDesc["sessionId"].(string)

	if !sessionId_ok {
		response["error"] = "sessionId required to logout"
		return c.JSON(http.StatusBadRequest, response)
	}

	AUTH_HANDLER.HandleLogout(sessionId)
	response["message"] = "logged out"
	return c.JSON(http.StatusOK, response)
}

func IsLoggedIn(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	sessionDesc, ok := body["data"].(map[string]any)
	response := map[string]any{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId, sessionId_ok := sessionDesc["sessionId"].(string)

	if !sessionId_ok {
		response["error"] = "sessionId required to check if user is logged in"
		return c.JSON(http.StatusBadRequest, response)
	}

	isLoggedIn := AUTH_HANDLER.IsLoggedIn(sessionId)
	response["isLoggedIn"] = isLoggedIn

	return c.JSON(http.StatusOK, response)
}

func ForgotPassword(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	email, ok := body["data"].(map[string]any)["email"].(string)
	if !ok {
		response := map[string]any{"error": "Invalid email format"}
		return c.JSON(http.StatusBadRequest, response)
	}

	passwordResetToken := AUTH_HANDLER.HandleForgotPassword(email)
	if passwordResetToken == "" {
		response := map[string]any{"error": "User not found or email sending failed"}
		return c.JSON(http.StatusNotFound, response)
	}

	resetLink := fmt.Sprintf("%s%s%s", config.BaseAuthUrl, "reset_password/", passwordResetToken)
	response := map[string]any{"message": "Password reset initiated. Please check your email for instructions.",
		"resetLink":          resetLink,
		"passwordResetToken": passwordResetToken}
	return c.JSON(http.StatusOK, response)
}

func ResetPassword(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	passwordResetToken := c.Param("passwordResetToken")

	newPassword, ok := body["data"].(map[string]any)["password"].(string)
	if !ok || newPassword == "" {
		response := map[string]any{"error": "Invalid or missing password"}
		return c.JSON(http.StatusBadRequest, response)
	}

	if !AUTH_HANDLER.HandleUpdatePassword(passwordResetToken, newPassword) {
		response := map[string]any{"error": "Failed to update password"}
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := map[string]any{"message": "Password updated successfully"}
	return c.JSON(http.StatusOK, response)
}
