package user_controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Iyusuf40/goBackendUtils/api/controllers"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
	"github.com/labstack/echo/v4"
)

var UserStorage = storage.MakeUserStorage(config.UsersDatabase,
	config.UsersRecords)
var SIGN_UP_HANDLER = controllers.MakeSignupHandler(config.TempStoreDb,
	config.UsersDatabase, config.UsersRecords)

func SaveUser(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	user := UserStorage.BuildClient(userDesc)

	if userWithEmailExist(user.Email) {
		response["error"] = "email already registered"
		return c.JSON(http.StatusBadRequest, response)
	}

	if config.RequireEmailVerification {
		signupId := SIGN_UP_HANDLER.HandleSignup(user)
		response["signupId"] = signupId
		return c.JSON(http.StatusCreated, response)
	}

	msg, success := UserStorage.Save(user)

	if !success {
		response["error"] = msg
		return c.JSON(http.StatusBadRequest, response)
	}

	response["userId"] = msg
	return c.JSON(http.StatusCreated, response)
}

func userWithEmailExist(email string) bool {
	users := UserStorage.GetByField("email", email)
	return len(users) >= 1
}

func CompleteSignup(c echo.Context) error {
	signupId := c.Param("signupId")
	response := map[string]string{}
	userId, success := SIGN_UP_HANDLER.HandleCompleteSignup(signupId)
	if !success {
		response["error"] = "failed to complete signup"
		return c.JSON(http.StatusBadRequest, response)
	}

	response["userId"] = userId
	return c.JSON(http.StatusCreated, response)
}

func GetUser(c echo.Context) error {
	userId := c.Param("id")
	user, err := UserStorage.Get(userId)
	response := map[string]string{}
	if err != nil {
		response["error"] = "user with id " + userId + " doesn't exist"
		return c.JSON(http.StatusNotFound, response)
	}

	return c.JSON(http.StatusOK, getUserMapWithoutPassword(user))
}

func UpdateUser(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	updateDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	field, fieldOk := updateDesc["field"].(string)
	value := updateDesc["value"]

	if !fieldOk {
		response["error"] = "field part of updateDesc is not a string"
		return c.JSON(http.StatusBadRequest, response)
	}

	userId := c.Param("id")

	updated := UserStorage.Update(userId, storage.UpdateDesc{Field: field,
		Value: value})

	if !updated {
		response["error"] = "update failed"
		return c.JSON(http.StatusBadRequest, response)
	}

	response["message"] = fmt.Sprintf("%s field of user succesfuly set to %s", field, value)
	return c.JSON(http.StatusOK, response)
}

func DeleteUser(c echo.Context) error {
	userId := c.Param("id")
	response := map[string]string{"message": "deleted"}

	UserStorage.Delete(userId)
	return c.JSON(http.StatusOK, response)
}

func getUserMapWithoutPassword(user models.User) map[string]any {
	userJSON, _ := json.Marshal(user)
	userMapRep := map[string]any{}
	json.Unmarshal(userJSON, &userMapRep)
	delete(userMapRep, "password")
	return userMapRep
}
