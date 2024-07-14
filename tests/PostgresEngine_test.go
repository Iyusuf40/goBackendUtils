package tests

import (
	"fmt"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/storage"
)

var database = "test"
var table = "users"
var POSTGRES_ENGINE *storage.PostgresEngine

func beforeEachPOSTGRES_ENGINE_T() {
	POSTGRES_ENGINE, _ = storage.MakePostgresEngine(database,
		table,
		storage.SQL_TABLE_COLUMN_FIELD_AND_DESC{"name", "varchar(256)"},
		storage.SQL_TABLE_COLUMN_FIELD_AND_DESC{"age", "integer"},
	)
}

func afterEachFPOSTGRES_ENGINE_T() {
	POSTGRES_ENGINE.DeleteTable()
	storage.RemovePostgressEngineSingleton(database, table, false)
}

func TestSaveAndGetPOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	user := User{"test user", 20}
	if id, err := POSTGRES_ENGINE.Save(user); err == nil {
		var obj, err = POSTGRES_ENGINE.Get(id)
		if obj == nil {
			t.Fatal("TestGet: early fail:", err.Error())
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

		// test get non-existing user
		obj, err = POSTGRES_ENGINE.Get("doNotExist")

		if obj != nil || err == nil {
			t.Fatal("TestGet: should fail geting non-existing user")
		}

	} else {
		t.Fatal("TestGet: failed to Save", err.Error())
	}
}

func TestGetRecordsByFieldPOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	user := User{"user", 20}
	POSTGRES_ENGINE.Save(user)
	var age float32 = 20

	records, err := POSTGRES_ENGINE.GetRecordsByField("age", age)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	nonExistentAge := 10
	records, err = POSTGRES_ENGINE.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}

	records, err = POSTGRES_ENGINE.GetRecordsByField("name", "user")

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	records, err = POSTGRES_ENGINE.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}
}

func TestGetIdByFieldAndFieldPOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	age := 20
	name := "userName"
	user := User{name, age}
	id, _ := POSTGRES_ENGINE.Save(user)

	retrievedId := POSTGRES_ENGINE.GetIdByFieldAndValue("age", age)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	retrievedId = POSTGRES_ENGINE.GetIdByFieldAndValue("name", name)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	// get non-existent user
	retrievedId = POSTGRES_ENGINE.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

	// get non-existent Record
	retrievedId = POSTGRES_ENGINE.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

}

func TestGetAllOfRecordsPOSTGRES_ENGINE(t *testing.T) {
	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	user := User{"user", 20}
	noOfSaves := 5

	for i := 0; i < noOfSaves; i++ {
		POSTGRES_ENGINE.Save(user)
	}
	allSavedUsers := POSTGRES_ENGINE.GetAllOfRecords()

	if len(allSavedUsers) != noOfSaves {
		t.Fatal("TestGetRecordsByField: length of records returned should be", noOfSaves)
	}
}

func TestUpdatePOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	user := User{"test", 20}
	id, _ := POSTGRES_ENGINE.Save(user)
	updated_name := "updated_name"

	resp := POSTGRES_ENGINE.Update(id, storage.UpdateDesc{
		Field: "name",
		Value: updated_name})

	if !resp {
		t.Fatal("TestUpdate: failed to update")
	}

	obj, _ := POSTGRES_ENGINE.Get(id)
	saved_user := new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}

	if POSTGRES_ENGINE.AllRecordsCount() != 1 {
		t.Fatal("TestUpdate: all records count should be 1")
	}
}

func TestDeletePOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	user := User{"test", 20}
	id, _ := POSTGRES_ENGINE.Save(user)

	if POSTGRES_ENGINE.AllRecordsCount() != 1 {
		t.Fatal("TestReload: records in db should be 1")
	}

	POSTGRES_ENGINE.Delete(id)

	if POSTGRES_ENGINE.AllRecordsCount() != 0 {
		t.Fatal("TestReload: records in db should be 0")
	}
}

func TestAllRecordsCountPOSTGRES_ENGINE(t *testing.T) {

	beforeEachPOSTGRES_ENGINE_T()
	defer afterEachFPOSTGRES_ENGINE_T()

	if POSTGRES_ENGINE.AllRecordsCount() != 0 {
		t.Fatal("TestAllRecordsCount: records in inMemoryStore should be 0")
	}

	noToSave := 10
	user := User{"test", 20}

	for i := 1; i <= noToSave; i++ {
		POSTGRES_ENGINE.Save(user)
		if POSTGRES_ENGINE.AllRecordsCount() != i {
			t.Fatal("TestAllRecordsCount: records in inMemoryStore should be " + fmt.Sprint(i))
		}
	}
}
