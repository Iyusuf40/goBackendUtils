package tests

import (
	"os"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/auth"
	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
)

var AuthHandler_userstorage_test_db_path = "test"
var AuthHandler_tempDBstorage_test_db_path = "test"
var Auth_Users_RecordsName = "users"

var AUTH_HANDLER *auth.AuthHandler
var AUTH_US storage.Storage[models.User]

func beforeEachAUTH_TEST() {
	AUTH_US = storage.MakeUserStorage(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
	AUTH_HANDLER = auth.MakeAuthHandler(AuthHandler_tempDBstorage_test_db_path,
		AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
}

func afterEachAUTH_TEST() {
	if config.DBMS == "postgres" {
		storage.RemovePostgressEngineSingleton(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName, true)
		storage.RemovePostgressEngineSingleton(AuthHandler_tempDBstorage_test_db_path, Auth_Users_RecordsName, true)
	} else if config.DBMS == "mongo" {
		storage.RemoveMongoSingleton(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName, true)
		storage.RemoveMongoSingleton(AuthHandler_tempDBstorage_test_db_path, Auth_Users_RecordsName, true)
	} else {
		storage.RemoveDbSingleton(AuthHandler_userstorage_test_db_path, Auth_Users_RecordsName)
		storage.RemoveDbSingleton(AuthHandler_tempDBstorage_test_db_path, Auth_Users_RecordsName)
		os.Remove(AuthHandler_userstorage_test_db_path)
		os.Remove(AuthHandler_tempDBstorage_test_db_path)
	}

}

func TestHandleLogin(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleLogin: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestHandleLogin: expected a sessionId got empty string")
	}

	// should fail to get a session Token
	sessId = AUTH_HANDLER.HandleLogin(email, password+"a")
	if sessId != "" {
		t.Fatal("TestHandleLogin: expected a sessionId to be empty string")
	}

	sessId = AUTH_HANDLER.HandleLogin(email+email, password)
	if sessId != "" {
		t.Fatal("TestHandleLogin: expected a sessionId to be empty string")
	}
}

func TestIsLoggedIn(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestIsLoggedIn: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestIsLoggedIn: expected a sessionId got empty string")
	}

	is_logged_in := AUTH_HANDLER.IsLoggedIn(sessId)

	if !is_logged_in {
		t.Fatal("TestIsLoggedIn: expected is_logged_in to be true")
	}

	is_logged_in = AUTH_HANDLER.IsLoggedIn(sessId + "abc")

	if is_logged_in {
		t.Fatal("TestIsLoggedIn: expected is_logged_in to be false")
	}
}

func TestHandleLogout(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleLogout: success should be true")
	}

	// should get a session Token
	sessId := AUTH_HANDLER.HandleLogin(email, password)
	if sessId == "" {
		t.Fatal("TestHandleLogout: expected a sessionId got empty string")
	}

	is_logged_in := AUTH_HANDLER.IsLoggedIn(sessId)

	if !is_logged_in {
		t.Fatal("TestHandleLogout: expected is_logged_in to be true")
	}

	// logout
	AUTH_HANDLER.HandleLogout(sessId)
	is_logged_in = AUTH_HANDLER.IsLoggedIn(sessId)

	if is_logged_in {
		t.Fatal("TestHandleLogout: expected is_logged_in to be false")
	}
}

func TestHandleForgotPassword(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	// test wrong email

	if AUTH_HANDLER.HandleForgotPassword("non-existent-email") != "" {
		t.Fatal("TestHandleForgotPassword: return value should be ''")
	}

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleForgotPassword: success should be true")
	}

	// test existing email
	passwordResetToken := AUTH_HANDLER.HandleForgotPassword(email)

	if passwordResetToken == "" {
		t.Fatal("TestHandleForgotPassword: passwordResetToken should not be empty")
	}
}

func TestHandleResetPassword(t *testing.T) {
	beforeEachAUTH_TEST()
	defer afterEachAUTH_TEST()

	// test wrong email

	newPassword := "new pass"

	if AUTH_HANDLER.HandleUpdatePassword("non-existent-email", newPassword) != false {
		t.Fatal("TestHandleUpdatePassword: return value should be false")
	}

	email := "testmail@mail.com"
	password := "xxx"

	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	_, success := AUTH_US.Save(user)

	if !success {
		t.Fatal("TestHandleUpdatePassword: success should be true")
	}

	passwordResetToken := AUTH_HANDLER.HandleForgotPassword(email)

	if passwordResetToken == "" {
		t.Fatal("TestHandleUpdatePassword: passwordResetToken should not be empty")
	}

	passwordUpdateIsSuccessful := AUTH_HANDLER.HandleUpdatePassword(passwordResetToken, newPassword)

	if !passwordUpdateIsSuccessful {
		t.Fatal("TestHandleUpdatePassword: passwordReset should be successful")
	}
}
