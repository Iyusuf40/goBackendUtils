package tests

import (
	"os"
	"slices"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/Iyusuf40/goBackendUtils/models"
	"github.com/Iyusuf40/goBackendUtils/storage"
)

var users_storage_test_db_path = "test"
var US storage.Storage[models.User]

func beforeEachUST() {
	US = storage.MakeUserStorage(users_storage_test_db_path, "users")
}

func afterEachUST() {
	if config.DBMS == "postgres" {
		storage.RemovePostgressEngineSingleton(users_storage_test_db_path, "users", true)
	} else if config.DBMS == "mongo" {
		storage.RemoveMongoSingleton(users_storage_test_db_path, "users", true)
	} else {
		storage.RemoveDbSingleton(users_storage_test_db_path, "users")
		os.Remove(users_storage_test_db_path)
	}

}

func TestSaveAndGetUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestSaveAndGetUser: success should be true;", id)
	}

	retrievedUser, _ := US.Get(id)

	if usersAreEqual(retrievedUser, user) == false {
		t.Fatal("TestSaveAndGetUser: retrievedUser should be equal to saved")
	}

	// test user with similar mail cannot be duplicated
	if _, success = US.Save(user); success {
		t.Fatal("TestSaveAndGetUser: success should be false")
	}
}

func TestUpdateUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	updateField := "phone"
	updateValue := 9000

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestUpdateUser: success should be true; msg:", id)
	}

	retrievedUser, _ := US.Get(id)

	if usersAreEqual(retrievedUser, user) == false {
		t.Fatal("TestUpdateUser: retrievedUser should be equal to saved")
	}

	US.Update(id, storage.UpdateDesc{Field: updateField, Value: updateValue})

	retrievedUser, _ = US.Get(id)

	if retrievedUser.Phone != updateValue {
		t.Fatal("TestUpdateUser: retrievedUser.Phone should be equal", updateValue)
	}
}

func TestUpdateUserPassword(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	password := "xxx"

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  password,
	}

	newPass := "yyy"

	updateField := "password"
	updateValue := newPass

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestUpdateUserPassword: success should be true; msg:", id)
	}

	retrievedUser, _ := US.Get(id)

	if usersAreEqual(retrievedUser, user) == false {
		t.Fatal("TestUpdateUserPassword: retrievedUser should be equal to saved")
	}

	US.Update(id, storage.UpdateDesc{Field: updateField, Value: updateValue})

	retrievedUser, _ = US.Get(id)

	if !retrievedUser.IsCorrectPassword(newPass) {
		t.Fatal("TestUpdateUser: retrievedUser.Phone should be equal", updateValue)
	}
}

func TestDeleteUser(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	user := models.User{
		Email:     "testmail@mail.com",
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestDeleteUser: success should be true")
	}

	retrievedUser, _ := US.Get(id)

	if usersAreEqual(retrievedUser, user) == false {
		t.Fatal("TestDeleteUser: retrievedUser should be equal to saved")
	}

	US.Delete(id)
	retrievedUser, err := US.Get(id)

	if usersAreEqual(retrievedUser, models.User{}) != true {
		t.Fatal("TestDeleteUser: retrievedUser should be empty")
	}

	if err == nil {
		t.Fatal("TestDeleteUser: getting nonexistent user should return error")
	}
}

func TestGetUserByField(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	email := "mail@mail"
	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	_, success := US.Save(user)
	if !success {
		t.Fatal("TestGetUserByField: success should be true")
	}

	retrievedUser := US.GetByField("email", email)[0]

	if usersAreEqual(retrievedUser, user) == false {
		t.Fatal("TestGetUserByField: retrievedUser should be equal to saved")
	}
}

func TestGetUserIdByField(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	email := "mail@mail"
	user := models.User{
		Email:     email,
		FirstName: "f_name",
		LastName:  "l_name",
		Phone:     8000,
		Password:  "xxx",
	}

	id, success := US.Save(user)
	if !success {
		t.Fatal("TestGetUserIdByField: success should be true")
	}

	retrievedId := US.GetIdByField("email", email)

	if retrievedId != id {
		t.Fatal("TestGetUserIdByField: retrievedId should be same")
	}

	// get wrong field
	retrievedId = US.GetIdByField("email", email+email)

	if retrievedId != "" {
		t.Fatal("TestGetUserIdByField: retrievedId should be empty")
	}
}

func TestGetAllUsers(t *testing.T) {
	beforeEachUST()
	defer afterEachUST()

	emails := []string{"mail1@mail.com", "mail2@mail.com", "mail3@mail.com"}

	for _, email := range emails {
		user := models.User{
			Email:     email,
			FirstName: "f_name",
			LastName:  "l_name",
			Phone:     8000,
			Password:  "xxx",
		}
		_, success := US.Save(user)
		if !success {
			t.Fatal("TestGetUserByField: success should be true")
		}
	}

	retrievedUsers := US.GetAll()

	if len(retrievedUsers) != len(emails) {
		t.Fatal("TestGetAllUsers: retrievedUsers should have length equal to saved users")
	}

	for _, user := range retrievedUsers {
		if slices.Contains(emails, user.Email) == false {
			t.Fatal("TestGetAllUsers:", user.Email, "should be in", emails)
		}
	}
}

func usersAreEqual(u1 models.User, u2 models.User) bool {
	if u1.Email != u2.Email ||
		u1.FirstName != u2.FirstName ||
		u1.LastName != u2.LastName ||
		u1.Phone != u2.Phone {
		return false
	}
	return true
}
