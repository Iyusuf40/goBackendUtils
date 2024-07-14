package tests

import (
	"fmt"
	"testing"

	"github.com/Iyusuf40/goBackendUtils/storage"
)

var mongo_database = "test"
var MONGO_WRAPPER *storage.MongoWrapper

func beforeEachMWRT() {
	MONGO_WRAPPER, _ = storage.MakeMongoWrapper(mongo_database, "User")
}

func afterEachMWRT() {
	storage.RemoveDbSingleton(test_db_path, "User")
	MONGO_WRAPPER.DeleteDb()
}

func TestSaveAndGetMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"test user", 20}
	if id, err := MONGO_WRAPPER.Save(user); err == nil {
		var obj, err = MONGO_WRAPPER.Get(id)
		fmt.Println(id, "<--------")
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

		// test non existing id
		obj, err = MONGO_WRAPPER.Get("non_existing")
		if obj != nil || err == nil {
			t.Fatal("TestGet: early fail")
		}

	} else {
		t.Fatal("TestGet: failed to Save")
	}
}

func TestGetRecordsByFieldMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"user", 20}
	MONGO_WRAPPER.Save(user)
	var age float32 = 20

	records, err := MONGO_WRAPPER.GetRecordsByField("age", age)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	nonExistentAge := 10
	records, err = MONGO_WRAPPER.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}

	records, err = MONGO_WRAPPER.GetRecordsByField("name", "user")

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 1 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 1",
			"got", len(records))
	}

	records, err = MONGO_WRAPPER.GetRecordsByField("age", nonExistentAge)

	if err != nil {
		t.Fatal(err)
	}

	if len(records) != 0 {
		t.Fatal("TestGetRecordsByField: length of records returned should be 0",
			"got", len(records))
	}
}

func TestGetIdByFieldAndFieldMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	age := 20
	name := "userName"
	user := User{name, age}
	id, _ := MONGO_WRAPPER.Save(user)

	retrievedId := MONGO_WRAPPER.GetIdByFieldAndValue("age", age)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	retrievedId = MONGO_WRAPPER.GetIdByFieldAndValue("name", name)

	if retrievedId != id {
		t.Fatal("TestGetIdByFieldAndField: expected ", id, "got", retrievedId)
	}

	// get non-existent user
	retrievedId = MONGO_WRAPPER.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

	// get non-existent Record
	retrievedId = MONGO_WRAPPER.GetIdByFieldAndValue("age", age+30)

	if retrievedId != "" {
		t.Fatal("TestGetIdByFieldAndField: expected ", "''", "got", retrievedId)
	}

}

func TestGetAllOfRecordsMWR(t *testing.T) {
	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"user", 20}
	noOfSaves := 5

	for i := 0; i < noOfSaves; i++ {
		MONGO_WRAPPER.Save(user)
	}
	allSavedUsers := MONGO_WRAPPER.GetAllOfRecords()

	if len(allSavedUsers) != noOfSaves {
		t.Fatal("TestGetRecordsByField: length of records returned should be", noOfSaves)
	}
}

func TestUpdateMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"test", 20}
	id, _ := MONGO_WRAPPER.Save(user)
	updated_name := "updated_name"

	resp := MONGO_WRAPPER.Update(id, storage.UpdateDesc{
		Field: "name",
		Value: updated_name})

	if !resp {
		t.Fatal("TestUpdate: failed to update")
	}

	obj, _ := MONGO_WRAPPER.Get(id)
	saved_user := new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}

	if MONGO_WRAPPER.AllRecordsCount() != 1 {
		t.Fatal("TestUpdate: all records count should be 1")
	}

	obj, _ = MONGO_WRAPPER.Get(id)
	saved_user = new(User).buildUser(obj)

	if saved_user.Name != updated_name {
		t.Fatal(
			"TestUpdate: failed to update Name field " +
				"expected " + updated_name + " got " + saved_user.Name,
		)
	}
}

func TestDeleteMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"test", 20}
	id, _ := MONGO_WRAPPER.Save(user)

	if MONGO_WRAPPER.AllRecordsCount() != 1 {
		t.Fatal("TestReload: records in inMemoryStore should be 1")
	}

	MONGO_WRAPPER.Delete(id)

	if MONGO_WRAPPER.AllRecordsCount() != 0 {
		t.Fatal("TestReload: records in inMemoryStore should be 0")
	}
}

func TestAllRecordsCountMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	if MONGO_WRAPPER.AllRecordsCount() != 0 {
		t.Fatal("TestAllRecordsCount: records in inMemoryStore should be 0")
	}

	noToSave := 10
	user := User{"test", 20}

	for i := 1; i <= noToSave; i++ {
		MONGO_WRAPPER.Save(user)
		if MONGO_WRAPPER.AllRecordsCount() != i {
			t.Fatal("TestAllRecordsCount: records in inMemoryStore should be " + fmt.Sprint(i))
		}
	}
}

func TestDeleteDbMWR(t *testing.T) {

	beforeEachMWRT()
	defer afterEachMWRT()

	user := User{"test", 20}

	MONGO_WRAPPER.Save(user)
	if MONGO_WRAPPER.AllRecordsCount() != 1 {
		t.Fatal("TestDeleteDbMWR: records in DB should be 1")
	}

	MONGO_WRAPPER.DeleteDb()

	if MONGO_WRAPPER.AllRecordsCount() != 0 {
		t.Fatal("TestDeleteDbMWR: records in DB should be 0, got")
	}
}
