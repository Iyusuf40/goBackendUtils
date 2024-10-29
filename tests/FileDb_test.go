package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/storage"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) buildUser(obj any) *User {
	if map_rep, ok := obj.(map[string]any); ok {
		usr_ob := new(User)
		usr_ob.Name = map_rep["name"].(string)
		age, _ := storage.GetFloat64Equivalent(map_rep["age"])
		usr_ob.Age = int(age)
		if usr_ob.Name == "" {
			return nil
		}
		return usr_ob
	}
	return nil
}

var test_db_path = "test_db.json"
var DB *storage.FileDb

func beforeEachFDBT() {
	DB, _ = storage.MakeFileDb(test_db_path, "User")
}

func afterEachFDBT() {
	storage.RemoveDbSingleton(test_db_path, "User")
	DB.DeleteDb()
}

func TestSaveAndGet(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	user := User{"test user", 20}
	if id, err := DB.Save(user); err == nil {
		var obj, _ = DB.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}

		saved_user := new(User).buildUser(obj)

		if saved_user.Name != user.Name {
			t.Fatal(
				"TestGet: name of user not equal: user.name = ",
				user.Name, "got =", saved_user.Name)
		}

		if saved_user.Age != user.Age {
			t.Fatal(
				"TestGet: age of user not equal: user.age =",
				user.Age, "got =", saved_user.Age)
		}

		// test load from file
		DB.Commit()
		DB.Reload()
		obj, _ = DB.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail")
		}

		saved_user = new(User).buildUser(obj)

		if saved_user.Name != user.Name {
			t.Fatal(
				"TestGet: name of user not equal: user.name = ",
				user.Name, "got =", saved_user.Name)
		}

		if saved_user.Age != user.Age {
			t.Fatal(
				"TestGet: age of user not equal: user.age =",
				user.Age, "got =", saved_user.Age)
		}

		DB.DeleteDb()

	} else {
		t.Fatal("TestGet: failed to Save")
	}
}

func TestGetRecordsByField(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	user := User{"user", 20}
	DB.Save(user)
	var age float32 = 20

	records, err := DB.GetRecordsByField("age", age)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	nonExistentAge := 10
	records, err = DB.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}

	records, err = DB.GetRecordsByField("name", "user")

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	records, err = DB.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}

	// ============================= test nested fields
	type Nested struct {
		A int
		B struct {
			BA int
		}
	}

	nstd := Nested{
		A: 1,
		B: struct {
			BA int
		}{1},
	}
	DB.Save(nstd)
	records, err = DB.GetRecordsByField("B.BA", 1)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}
}

func TestGetIdByFieldAndField(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	age := 20
	name := "userName"
	user := User{name, age}
	id, _ := DB.Save(user)

	retrievedId := DB.GetIdByFieldAndValue("age", age)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	retrievedId = DB.GetIdByFieldAndValue("name", name)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	// get non-existent user
	retrievedId = DB.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

	// get non-existent Record
	retrievedId = DB.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

}

func TestGetAllOfRecords(t *testing.T) {
	beforeEachFDBT()
	defer afterEachFDBT()

	user := User{"user", 20}
	noOfSaves := 5

	for i := 0; i < noOfSaves; i++ {
		DB.Save(user)
	}
	allSavedUsers := DB.GetAllOfRecords()

	if len(allSavedUsers) != noOfSaves {
		t.Fatal("TestGetRecordsByField: length of records returned should be", noOfSaves)
	}
}

func TestUpdate(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	user := User{"test", 20}
	id, _ := DB.Save(user)
	updated_name := "updated_name"

	resp := DB.Update(id, storage.UpdateDesc{
		Field: "name",
		Value: updated_name})

	if !resp {
		t.Fatal("TestUpdate: failed to update")
	}

	obj, _ := DB.Get(id)
	saved_user := new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}

	// test after reload
	DB.Commit()
	DB.Reload()

	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestUpdate: all records count should be 1")
	}

	obj, _ = DB.Get(id)
	saved_user = new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}

	// ============================= test nested fields
	type Nested struct {
		A int
		B struct {
			BA int
		}
	}

	nstd := Nested{
		A: 1,
		B: struct {
			BA int
		}{1},
	}

	id, _ = DB.Save(nstd)

	updated_nested_field := 2

	resp = DB.Update(id, storage.UpdateDesc{
		Field: "B.BA",
		Value: updated_nested_field})

	if !resp {
		t.Fatal("TestUpdate: failed to update")
	}

	obj, _ = DB.Get(id)

	if obj.(map[string]any)["B"].(map[string]any)["BA"] != 2 {
		t.Fatal(
			"TestUpdate: failed to update nested field, expected",
			2, "got", obj.(map[string]any)["B"].(map[string]any)["BA"])
	}
}

func TestDelete(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	user := User{"test", 20}
	id, _ := DB.Save(user)

	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestReload: records in inMemoryStore should be 1")
	}

	DB.Delete(id)

	if DB.AllRecordsCount() != 0 {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}
}

func TestAllRecordsCount(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	if DB.AllRecordsCount() != 0 {
		t.Fatal("TestAllRecordsCount: records in inMemoryStore should be 0")
	}

	noToSave := 10
	user := User{"test", 20}

	for i := 1; i <= noToSave; i++ {
		DB.Save(user)
		if DB.AllRecordsCount() != i {
			t.Fatal("TestAllRecordsCount: records in inMemoryStore should be " + fmt.Sprint(i))
		}
	}
}

func TestReload(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	zero := 0
	// nothing in inMemoryStore
	if DB.AllRecordsCount() != zero {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}

	user := User{"test", 20}
	DB.Save(user)

	// Reload should only load committed transactions
	DB.Reload()
	if DB.AllRecordsCount() != zero {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}

	id, _ := DB.Save(user)
	DB.Commit()
	DB.Reload()
	if DB.AllRecordsCount() != 1 {
		t.Fatal("TestReload: 1 record should be in inMemoryStore")
	}

	got, _ := DB.Get(id)
	saved_user := new(User).buildUser(got)

	if saved_user.Name != user.Name {
		t.Fatal(
			"TestReload: name of user not equal: user.name = ",
			user.Name, "got =", saved_user.Name)
	}

	if saved_user.Age != user.Age {
		t.Fatal(
			"TestReload: age of user not equal: user.age =",
			user.Age, "got =", saved_user.Age)
	}

}

func TestCommit_DeleteDb(t *testing.T) {

	beforeEachFDBT()
	defer afterEachFDBT()

	// test db_file does not exist
	_, err := os.Stat(test_db_path)
	if err == nil {
		t.Fatal("TestCommit_DeleteDb: db_file should not exist")
	}

	// test db_file should exist after commit
	DB.Commit()
	_, err = os.Stat(test_db_path)
	if err != nil {
		t.Fatal("TestCommit_DeleteDb: db_file should exist")
	}

	// test db_file should not exist after DeleteDb
	DB.DeleteDb()
	_, err = os.Stat(test_db_path)
	if err == nil {
		t.Fatal("TestCommit_DeleteDb: db_file should not exist")
	}
}
