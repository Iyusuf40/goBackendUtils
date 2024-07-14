package storage

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
)

type FileDb struct {
	path                       string
	recordsName                string
	inMemoryStore              map[string]any
	RECORDS_NAME_KEY_SEPARATOR string
}

func (db *FileDb) New(db_path, recordsName string) (*FileDb, error) {
	if db_path == "" || recordsName == "" {
		panic("FileDb.New: db_path and objectType must not be empty")
	}
	db.path = db_path
	db.recordsName = recordsName
	db.RECORDS_NAME_KEY_SEPARATOR = "-"
	err := db.Reload()
	return db, err
}

func (db *FileDb) Reload() error {
	db.inMemoryStore = make(map[string]any)
	content, _ := os.ReadFile(db.path)

	json.Unmarshal(content, &(db.inMemoryStore))

	return nil
}

func (db *FileDb) AllRecordsCount() int {
	return len(db.inMemoryStore)
}

func (db *FileDb) Save(obj any) (string, error) {
	uid := uuid.NewString()
	id := db.recordsName + db.RECORDS_NAME_KEY_SEPARATOR + uid
	json_rep, err := json.Marshal(obj) // test if it can be jsoned
	if err != nil {
		return "", err
	}

	// save map[string]any rep
	var saved_version map[string]any
	json.Unmarshal(json_rep, &saved_version)

	db.inMemoryStore[id] = saved_version

	return id, nil
}

// returns objects with any type so users can rebuild
// objects with their type builders
func (db *FileDb) Get(id string) (any, error) {
	stored, found := db.inMemoryStore[id]
	if found {
		return stored, nil
	}
	return nil, errors.New("FileDb: Get: failed to get object with id: " + id)
}

func (db *FileDb) GetRecordsByField(field string, value any) ([]map[string]any, error) {
	var listOfRecordsOfSameType = db.GetAllOfRecords()

	var listOfMatchedRecords []map[string]any
	var compValue any

	// convert value number to float64 if value is a number
	numberVal, ok := getFloat64Equivalent(value)
	if ok {
		compValue = numberVal
	} else {
		compValue = value
	}

	for _, record := range listOfRecordsOfSameType {
		if record[field] == compValue {
			listOfMatchedRecords = append(listOfMatchedRecords, record)
		}
	}

	return listOfMatchedRecords, nil
}

func (db *FileDb) GetIdByFieldAndValue(field string, value any) string {

	recordsName := db.recordsName
	for key, val := range db.inMemoryStore {
		if strings.HasPrefix(key, recordsName) {
			concVal, ok := val.(map[string]any)
			if !ok {
				panic(`FileDb: GetRecordsByField: records found for is not of  
					map[string]any type` + recordsName)
			}
			if concVal[field] == value {
				return key
			}

			if floatRep, ok := concVal[field].(float64); ok {
				valueFloat, _ := getFloat64Equivalent(value)
				if floatRep == valueFloat {
					return key
				}
			}
		}
	}
	return ""
}

func (db *FileDb) GetAllOfRecords() []map[string]any {
	var listOfRecordsOfSameType []map[string]any
	recordsName := db.recordsName
	for key, val := range db.inMemoryStore {
		if strings.HasPrefix(key, recordsName) {
			concVal, ok := val.(map[string]any)
			if !ok {
				panic(`FileDb: GetRecordsByField: records found for is not of  
					map[string]any type` + recordsName)
			}
			listOfRecordsOfSameType = append(listOfRecordsOfSameType, concVal)
		}
	}

	return listOfRecordsOfSameType
}

func (db *FileDb) Delete(id string) {
	delete(db.inMemoryStore, id)
}

func (db *FileDb) Update(id string, data UpdateDesc) bool {
	_, exists := db.inMemoryStore[id]
	if !exists {
		return false
	}

	obj := db.inMemoryStore[id]

	if concrete_obj, ok := obj.(map[string]any); ok {
		concrete_obj[data.Field] = data.Value
		db.inMemoryStore[id] = concrete_obj
	} else {
		panic("typeof inMemoryStore[id] is not map[string]any")
	}

	return true
}

func (db *FileDb) Commit() error {
	json_rep, err := json.Marshal(db.inMemoryStore)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, json_rep, 0644)
	return err
}

func (db *FileDb) DeleteDb() error {
	delete(FILE_DB_MAP, db.path)
	err := os.Remove(db.path)
	return err
}

func GetFloat64Equivalent(value any) (float64, bool) {
	return getFloat64Equivalent(value)
}

func getFloat64Equivalent(value any) (float64, bool) {
	if concVal, ok := value.(int); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int8); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int16); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int32); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(int64); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(float32); ok {
		return float64(concVal), true
	}

	if concVal, ok := value.(float64); ok {
		return float64(concVal), true
	}

	return 0, false
}

var FILE_DB_MAP = map[string]*FileDb{}

func MakeFileDb(db_path string, recordsName string) (*FileDb, error) {
	path := db_path

	if recordsName == "" {
		panic("MakeFileDb: recordsName cannot be empty")
	}

	if path == "" {
		panic("MakeFileDb: db_path cannot be empty")
	}

	key := path + recordsName
	// implements singleton pattern
	if FILE_DB_MAP[key] != nil {
		return FILE_DB_MAP[key], nil
	}

	file_db, err := new(FileDb).New(path, recordsName)

	if err != nil {
		panic("MakeFileDb: " + err.Error())
	}

	FILE_DB_MAP[key] = file_db
	return file_db, nil
}

func RemoveDbSingleton(db_path, recordsName string) {
	if recordsName == "" {
		panic("RemoveDbSingleton: recordsName cannot be empty")
	}

	if db_path == "" {
		panic("RemoveDbSingleton: db_path cannot be empty")
	}

	key := db_path + recordsName
	delete(FILE_DB_MAP, key)
}
