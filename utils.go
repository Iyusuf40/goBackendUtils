package goBackendUtils

import (
	"io"

	"github.com/Iyusuf40/goBackendUtils/api/controllers"
	"github.com/Iyusuf40/goBackendUtils/storage"
	"github.com/labstack/echo/v4"
)

type Utils struct {
}

func (ut *Utils) GetBodyInMap(c echo.Context) map[string]any {
	return controllers.GetBodyInMap(c)
}

func (ut *Utils) ReadFromReaderIntoMap(r io.Reader) map[string]any {
	return controllers.ReadFromReaderIntoMap(r)
}

func (ut *Utils) GetDB_Engine(engine_dbms, database, recordsName string, fieldAndDesc ...storage.SQL_TABLE_COLUMN_FIELD_AND_DESC) (storage.DB_Engine, error) {
	return storage.GetDB_Engine(engine_dbms, database, recordsName, fieldAndDesc...)
}

func (ut *Utils) GET_TempStore(typ, database, recordsName string) storage.TempStore {
	return storage.GET_TempStore(typ, database, recordsName)
}
