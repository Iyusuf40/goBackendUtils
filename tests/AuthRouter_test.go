package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/api/controllers"
	"github.com/Iyusuf40/goBackendUtils/auth"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
	"github.com/labstack/echo/v4"
)

var auth_test_db_path = "test"
var auth_test_user_recordsName = "users"
var USER_STORE storage.Storage[models.User]

func beforeEachUAUTHT() {
	USER_STORE = storage.MakeUserStorage(auth_test_db_path, auth_test_user_recordsName)
	auth.AUTH_HANDLER = auth.MakeAuthHandler(auth_test_db_path, auth_test_db_path, auth_test_user_recordsName)
}

func afterEachUAUTHT() {
	if config.DBMS == "postgres" {
		storage.RemovePostgressEngineSingleton(auth_test_db_path, auth_test_user_recordsName, true)
	} else if config.DBMS == "mongo" {
		storage.RemoveMongoSingleton(auth_test_db_path, auth_test_user_recordsName, true)
	} else {
		storage.RemoveDbSingleton(auth_test_db_path, auth_test_user_recordsName)
		os.Remove(auth_test_db_path)
	}
}

func TestLoginUser(t *testing.T) {
	// Setup
	beforeEachUAUTHT()
	defer afterEachUAUTHT()

	e := echo.New()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	_, success := USER_STORE.Save(user)

	if !success {
		t.Fatal("TestLoginUser: success should be true")
	}

	loginDataJSON := `{"data": {"email":"testmail@mail.com", "password": "xxx"}}`

	// test successfully login a user
	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}

	rec, c := SetupRequest(e, http.MethodPost, "/auth/login", loginDataJSON, headers)
	auth.Login(c)

	if rec.Code != http.StatusOK {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: expected:", http.StatusOK, "got:", rec.Code)
	}

	recBody := controllers.ReadFromReaderIntoMap(rec.Body)

	if sessId, ok := recBody["sessionId"].(string); !ok || sessId == "" {
		t.Fatal("POST /auth/login: expected sessionId to exist")
	}

	sessId := recBody["sessionId"].(string)

	// test isLoggedIn
	isloggedinDataJSON := fmt.Sprintf(`{"data": {"sessionId":"%s"}}`, sessId)
	rec, c = SetupRequest(e, http.MethodPost, "/auth/isloggedin", isloggedinDataJSON, headers)
	auth.IsLoggedIn(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: expected:", http.StatusOK, "got:", rec.Code)
	}

	// test failed login of user with wrong password
	loginDataJSON = `{"data": {"email":"testmail@mail.com", "password": "yyy"}}`

	rec, c = SetupRequest(e, http.MethodPost, "/api/users", loginDataJSON, headers)
	auth.Login(c)

	if http.StatusBadRequest != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: expected:", http.StatusBadRequest, "got:", rec.Code)
	}

	recBody = controllers.ReadFromReaderIntoMap(rec.Body)

	if sessId, ok := recBody["sessionId"].(string); ok || sessId != "" {
		t.Fatal("POST /auth/login: expected sessionId to not exist")
	}

	// test failed login of user with missing required fields
	loginDataJSON = `{"data": {"em":"", "password": ""}}`

	rec, c = SetupRequest(e, http.MethodPost, "/api/users", loginDataJSON, headers)
	auth.Login(c)

	if http.StatusBadRequest != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: expected:", http.StatusBadRequest, "got:", rec.Code)
	}

	// test failed login of user with missing required fields
	loginDataJSON = `{"data": {"email":"", "password": ""}}`

	rec, c = SetupRequest(e, http.MethodPost, "/auth/login", loginDataJSON, headers)
	auth.Login(c)

	if http.StatusBadRequest != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: expected:", http.StatusBadRequest, "got:", rec.Code)
	}
}

func TestLogoutUser(t *testing.T) {
	// Setup
	beforeEachUAUTHT()
	defer afterEachUAUTHT()

	e := echo.New()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	_, success := USER_STORE.Save(user)

	if !success {
		t.Fatal("TestLogoutUser: success should be true")
	}

	loginDataJSON := `{"data": {"email":"testmail@mail.com", "password": "xxx"}}`

	// test successfully login a user
	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}

	rec, c := SetupRequest(e, http.MethodPost, "/auth/login", loginDataJSON, headers)
	auth.Login(c)

	// login user
	if rec.Code != http.StatusOK {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/logout: expected:", http.StatusOK, "got:", rec.Code)
	}

	recBody := controllers.ReadFromReaderIntoMap(rec.Body)

	if sessId, ok := recBody["sessionId"].(string); !ok || sessId == "" {
		t.Fatal("POST /auth/logout: expected sessionId to exist")
	}

	sessId := recBody["sessionId"].(string)

	// test isLoggedIn
	isloggedinDataJSON := fmt.Sprintf(`{"data": {"sessionId":"%s"}}`, sessId)
	rec, c = SetupRequest(e, http.MethodPost, "/auth/isloggedin", isloggedinDataJSON, headers)
	auth.IsLoggedIn(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/logout: expected:", http.StatusOK, "got:", rec.Code)
	}

	isLoggedIn := controllers.ReadFromReaderIntoMap(rec.Body)["isLoggedIn"].(bool)

	if !isLoggedIn {
		t.Fatal("POST /auth/logout: expected isLoggedIn to be true", "got:", isLoggedIn)
	}

	rec, c = SetupRequest(e, http.MethodPost, "/auth/logout", isloggedinDataJSON, headers)
	auth.Logout(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/logout: expected:", http.StatusOK, "got:", rec.Code)
	}

	// test isLoggedIn
	isloggedinDataJSON = fmt.Sprintf(`{"data": {"sessionId":"%s"}}`, sessId)
	rec, c = SetupRequest(e, http.MethodPost, "/auth/isloggedin", isloggedinDataJSON, headers)
	auth.IsLoggedIn(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/logout: expected:", http.StatusOK, "got:", rec.Code)
	}

	isLoggedIn = controllers.ReadFromReaderIntoMap(rec.Body)["isLoggedIn"].(bool)

	if isLoggedIn {
		t.Fatal("POST /auth/logout: expected isLoggedIn to be false", "got:", isLoggedIn)
	}

}

func TestForgotPassword(t *testing.T) {
	// Setup
	beforeEachUAUTHT()
	defer afterEachUAUTHT()

	e := echo.New()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	_, success := USER_STORE.Save(user)

	if !success {
		t.Fatal("TestForgotPassword: success should be true")
	}

	// test wrong email
	forgotPasswordJSON := `{"data": {"email":"wrongmail@mail.com"}}`

	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}

	rec, c := SetupRequest(e, http.MethodPost, "/auth/fotgot_password", forgotPasswordJSON, headers)
	auth.ForgotPassword(c)

	if http.StatusNotFound != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/forgot_password: expected:", http.StatusNotFound, "got:", rec.Code)
	}

	forgotPasswordJSON = `{"data": {"email":"testmail@mail.com"}}`

	rec, c = SetupRequest(e, http.MethodPost, "/auth/fotgot_password", forgotPasswordJSON, headers)
	auth.ForgotPassword(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/forgot_password: expected:", http.StatusOK, "got:", rec.Code)
	}
}

func TestResetPassword(t *testing.T) {
	// Setup
	beforeEachUAUTHT()
	defer afterEachUAUTHT()

	e := echo.New()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	_, success := USER_STORE.Save(user)

	if !success {
		t.Fatal("TestForgotPassword: success should be true")
	}

	headers := map[string]string{
		echo.HeaderContentType: echo.MIMEApplicationJSON,
	}

	newPassword := "newPassword"

	loginDataJSON := fmt.Sprintf(`{"data": {"email":"testmail@mail.com", "password": "%s"}}`, newPassword)

	// test login with newPassword before update password
	rec, c := SetupRequest(e, http.MethodPost, "/auth/login", loginDataJSON, headers)
	auth.Login(c)

	// login user should fail
	if rec.Code != http.StatusBadRequest {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: in TestResetPassword expected:", http.StatusBadRequest, "got:", rec.Code)
	}

	// initiate forgot_password process
	forgotPasswordJSON := `{"data": {"email":"testmail@mail.com"}}`

	rec, c = SetupRequest(e, http.MethodPost, "/auth/fotgot_password", forgotPasswordJSON, headers)
	auth.ForgotPassword(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/forgot_password: in TestResetPassword: expected:", http.StatusOK, "got:", rec.Code)
	}

	recBody := controllers.ReadFromReaderIntoMap(rec.Body)

	passwordResetToken := recBody["passwordResetToken"].(string)
	newPasswordJSON := `{"data": {"password":"newPassword"}}`

	rec, c = SetupRequest(e, http.MethodPost, "/auth/reset_password/"+passwordResetToken,
		newPasswordJSON, headers)

	c.SetParamNames("passwordResetToken")
	c.SetParamValues(passwordResetToken)
	auth.ResetPassword(c)

	if http.StatusOK != rec.Code {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/reset_password: expected:", http.StatusOK, "got:", rec.Code)
	}

	// test login with newPassword after update password
	rec, c = SetupRequest(e, http.MethodPost, "/auth/login", loginDataJSON, headers)
	auth.Login(c)

	// login user should fail
	if rec.Code != http.StatusOK {
		fmt.Println("body returned", rec.Body.String())
		t.Fatal("POST /auth/login: in TestResetPassword expected:", http.StatusOK, "got:", rec.Code)
	}
}
