package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/Iyusuf40/goBackendUtils/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresEngine struct {
	tableName string
	conn      *pgx.Conn
}

type SQL_TABLE_COLUMN_FIELD_AND_DESC [2]string

func (db *PostgresEngine) New(database, tableName string, fieldAndDesc ...SQL_TABLE_COLUMN_FIELD_AND_DESC) (*PostgresEngine, error) {
	if database == "" || tableName == "" {
		panic("PostgresEngine.New: db_path and objectType must not be empty")
	}
	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		config.DB_USER, config.DB_PASSWORD, config.DB_HOST, database)
	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PostgresEngine.New: Unable to connect to database: %v", err)
		return nil, err
	}

	db.tableName = tableName

	createTableStmt := db.makeCreateTableStmt(fieldAndDesc...)

	_, err = conn.Exec(context.Background(), createTableStmt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "PostgresEngine.New: Failed to create table: %v", err)
		return nil, err
	}

	db.conn = conn

	return db, err
}

// makeCreateTableStmt: creates the statement to create the table.
// params - fieldAndDesc ([2]string): the field is retrieved from index 0
// while the type, constaraint and all other field description is
// retrieved from the second index
func (db *PostgresEngine) makeCreateTableStmt(fieldAndDesc ...SQL_TABLE_COLUMN_FIELD_AND_DESC) string {

	// create implicit id column
	fieldAndDesc = append(fieldAndDesc,
		SQL_TABLE_COLUMN_FIELD_AND_DESC{"id", "varchar(64) PRIMARY KEY"})

	// make sure tables columns are created in a sorted manner
	// hence during insertion we just need to sort the object's
	// fields and it will match table creation column order
	sort.Slice(fieldAndDesc, func(i, j int) bool {
		return fieldAndDesc[i][0] < fieldAndDesc[j][0]
	})

	stmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (`, db.tableName)

	for i, fieldAndType := range fieldAndDesc {
		field, description := fieldAndType[0], fieldAndType[1]
		if i != len(fieldAndDesc)-1 {
			stmt += fmt.Sprintf(`"%s"		%s,`, field, description)
		} else {
			stmt += fmt.Sprintf(`"%s"		%s`, field, description)
		}
	}

	stmt += `)`

	return stmt
}

func (db *PostgresEngine) AllRecordsCount() int {
	stmt := fmt.Sprintf(`SELECT COUNT(*) FROM "%s";`, db.tableName)
	var count int
	db.conn.QueryRow(context.Background(), stmt).Scan(&count)
	return count
}

func (db *PostgresEngine) Save(obj any) (string, error) {
	id := uuid.NewString()
	json_rep, err := json.Marshal(obj) // test if it can be jsoned
	if err != nil {
		return "", err
	}

	var mapRep map[string]any
	json.Unmarshal(json_rep, &mapRep)

	mapRep["id"] = id

	insertStmt, parameters := db.makeInsertStmtAndParameters(mapRep)

	_, err = db.conn.Exec(context.Background(), insertStmt, parameters...)

	if err != nil {
		return "", err
	}

	return id, nil
}

// makeInsertStmtAndParameters - creates an insert statment from mapRep.
// param mapRep - a map of all the columns and their values.
// makeInsertStmt constructs the insert statment by sorting
// the columns alphabetically, this assumes the create table
// statement did the same during creation.
// returns - the function returns both statement with its positional
// params ($pos) embedded as well as the parameters
func (db *PostgresEngine) makeInsertStmtAndParameters(mapRep map[string]any) (string, []any) {

	fieldsAndValues := [][2]any{}

	for field, value := range mapRep {
		fieldsAndValues = append(fieldsAndValues, [2]any{field, value})
	}

	sort.Slice(fieldsAndValues, func(i, j int) bool {
		field1 := fieldsAndValues[i][0].(string)
		field2 := fieldsAndValues[j][0].(string)
		return field1 < field2
	})

	stmt := fmt.Sprintf(`INSERT INTO "%s"`, db.tableName)

	values := []any{}
	fields := `(`
	// use placeholder positional params to be substituted in the prepared stmt
	valuesPlaceHolder := `(`
	for index, fieldAndValue := range fieldsAndValues {
		field := fieldAndValue[0].(string)
		value := fieldAndValue[1]
		// append values in order of fields
		// to be returned with the parametarized statement for insertion
		values = append(values, value)
		if index != len(fieldsAndValues)-1 {
			fields += fmt.Sprintf(`"%s",`, field)
			valuesPlaceHolder += fmt.Sprintf(`$%v,`, index+1)
		} else {
			fields += fmt.Sprintf(`"%s"`, field)
			valuesPlaceHolder += fmt.Sprintf(`$%v`, index+1)
		}
	}
	fields += `) VALUES`
	valuesPlaceHolder += `);`

	stmt = fmt.Sprintf(`%s %s %s`, stmt, fields, valuesPlaceHolder)
	return stmt, values
}

// returns objects with any type so users can rebuild
// objects with their type builders
func (db *PostgresEngine) Get(id string) (any, error) {
	stmt := fmt.Sprintf(`SELECT * FROM "%s" WHERE id = $1;`, db.tableName)
	row, err := db.conn.Query(context.Background(), stmt, id)

	if err != nil {
		return nil, err
	}

	mapRep, err := pgx.CollectOneRow(row, pgx.RowToMap)
	if err != nil {
		return nil, err
	}
	return mapRep, nil
}

func (db *PostgresEngine) GetRecordsByField(field string, value any) ([]map[string]any, error) {
	stmt := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s" = $1;`, db.tableName, field)

	row, err := db.conn.Query(context.Background(), stmt, value)

	if err != nil {
		return nil, err
	}

	listOfmapReps, err := pgx.CollectRows(row, pgx.RowToMap)
	if err != nil {
		return nil, err
	}
	return listOfmapReps, nil
}

func (db *PostgresEngine) GetIdByFieldAndValue(field string, value any) string {
	listOfmapReps, err := db.GetRecordsByField(field, value)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return ""
	}

	if len(listOfmapReps) > 1 {
		err = errors.New("PostgresEngine.GetIdByFieldAndValue: returned list cannot be more than 1")
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if len(listOfmapReps) == 1 {
		return listOfmapReps[0]["id"].(string)
	}

	return ""
}

func (db *PostgresEngine) GetAllOfRecords() []map[string]any {
	stmt := fmt.Sprintf(`SELECT * FROM "%s";`, db.tableName)

	row, _ := db.conn.Query(context.Background(), stmt)

	listOfmapReps, _ := pgx.CollectRows(row, pgx.RowToMap)

	return listOfmapReps
}

func (db *PostgresEngine) Delete(id string) {
	stmt := fmt.Sprintf(`DELETE FROM "%s" WHERE id = $1;`, db.tableName)
	_, err := db.conn.Exec(context.Background(), stmt, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func (db *PostgresEngine) Update(id string, data UpdateDesc) bool {
	stmt := fmt.Sprintf(`UPDATE "%s" SET "%s" = $1 WHERE id = $2;`, db.tableName, data.Field)
	cmdTag, err := db.conn.Exec(context.Background(), stmt, data.Value, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	return cmdTag.RowsAffected() == 1
}

func (db *PostgresEngine) Commit() error {
	return nil
}

func (db *PostgresEngine) CloseConnection() error {
	return db.conn.Close(context.Background())
}

func (db *PostgresEngine) DeleteTable() error {
	_, err := db.conn.Exec(context.Background(), fmt.Sprintf("DROP TABLE %s;", db.tableName))
	return err
}

var POSTGRES_ENGINE_MAP = map[string]*PostgresEngine{}

func MakePostgresEngine(database string, tableName string, fieldAndDesc ...SQL_TABLE_COLUMN_FIELD_AND_DESC) (*PostgresEngine, error) {
	if tableName == "" {
		panic("MakePostgresEngine: tableName cannot be empty")
	}

	if database == "" {
		panic("MakePostgresEngine: database cannot be empty")
	}

	key := database + tableName

	// implements singleton pattern
	if POSTGRES_ENGINE_MAP[key] != nil {
		return POSTGRES_ENGINE_MAP[key], nil
	}

	postgresEng, err := new(PostgresEngine).New(database, tableName, fieldAndDesc...)

	if err != nil {
		panic("MakePostgresEngine: " + err.Error())
	}

	POSTGRES_ENGINE_MAP[key] = postgresEng
	return postgresEng, nil
}

func RemovePostgressEngineSingleton(database, tableName string, shouldDeleteTable ...bool) {
	if tableName == "" {
		panic("RemoveDbSingleton: tableName cannot be empty")
	}

	if database == "" {
		panic("RemoveDbSingleton: database cannot be empty")
	}

	key := database + tableName
	postgresEng, exists := POSTGRES_ENGINE_MAP[key]
	if exists {
		// this is just to make deleting a table dificult and intentional
		if len(shouldDeleteTable) == 1 {
			if shouldDeleteTable[0] && database == "test" {
				postgresEng.DeleteTable()
			}
		}

		if len(shouldDeleteTable) == 2 {
			if shouldDeleteTable[0] && shouldDeleteTable[1] {
				postgresEng.DeleteTable()
			}
		}

		postgresEng.CloseConnection()
		delete(POSTGRES_ENGINE_MAP, key)
	}
}
