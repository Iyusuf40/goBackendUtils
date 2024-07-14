package storage

type UpdateDesc struct {
	Field string
	Value any
}

// Every table or collection must implement the Storage interface
type Storage[T any] interface {
	Get(id string) (T, error)
	Save(data T) (msg string, success bool)
	Update(id string, data UpdateDesc) bool
	Delete(id string)
	GetByField(field string, value any) []T
	GetIdByField(field string, value any) string
	GetAll() []T
	BuildClient(obj any) T
}

type DB_Engine interface {
	Get(id string) (any, error)
	Save(data any) (string, error)
	// for document stores using noSql, Callers of this function
	// must validate both fields passed in data else, unwanted fields
	// may be added to the records on disc and values of
	// an inappropriate type might be added, causing errors in
	// rebuilding objects
	Update(id string, data UpdateDesc) bool
	Delete(id string)
	// if FileDb is the Engine, field is the json tag if it
	// is defined on the obj
	GetRecordsByField(field string, value any) ([]map[string]any, error)
	GetIdByFieldAndValue(field string, value any) string
	GetAllOfRecords() []map[string]any
	Commit() error
}

func GetDB_Engine(engine_dbms, database, recordsName string, fieldAndDesc ...SQL_TABLE_COLUMN_FIELD_AND_DESC) (DB_Engine, error) {
	switch engine_dbms {
	case "postgres":
		return MakePostgresEngine(database, recordsName, fieldAndDesc...)
	case "mongo", "mongodb":
		return MakeMongoWrapper(database, recordsName)
	default:
		return MakeFileDb(database, recordsName)
	}
}
