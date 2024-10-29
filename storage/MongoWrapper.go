package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoWrapper struct {
	client        *mongo.Client
	collection    *mongo.Collection
	database_name string
}

func (db *MongoWrapper) New(database, collection string) (*MongoWrapper, error) {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db.client = client
	db.database_name = database
	db.collection = client.Database(database).Collection(collection)
	return db, err
}

func (db *MongoWrapper) AllRecordsCount() int {
	count, err := db.collection.CountDocuments(context.TODO(), bson.D{})

	if err != nil {
		return -1
	}

	return int(count)
}

func (db *MongoWrapper) Save(obj any) (string, error) {

	id := ""

	json_rep, err := json.Marshal(obj) // test if it can be jsoned
	if err != nil {
		return "", err
	}

	var mapRep map[string]any
	json.Unmarshal(json_rep, &mapRep)

	bsonD := db.makeBsonDSlice(mapRep)
	result, err := db.collection.InsertOne(context.Background(), bsonD)

	if err != nil {
		return "", err
	}

	id = result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (db *MongoWrapper) makeBsonDSlice(mapRep map[string]any) bson.D {
	bsonD := bson.D{}

	for key, value := range mapRep {
		bsonD = append(bsonD, bson.E{Key: key, Value: value})
	}

	return bsonD
}

// returns objects with any type so users can rebuild
// objects with their type builders
func (db *MongoWrapper) Get(id string) (any, error) {
	var result any
	objectId, _ := primitive.ObjectIDFromHex(id)
	err := db.collection.FindOne(context.Background(),
		bson.D{{Key: "_id", Value: objectId}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	jsonRep, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	var sliceRep = []map[string]any{}

	err = json.Unmarshal(jsonRep, &sliceRep)

	if err != nil {
		return nil, err
	}

	objectAsMap := db.getObjectAsMap(sliceRep)
	return objectAsMap, nil
}

func (db *MongoWrapper) getObjectAsMap(sliceRep []map[string]any) map[string]any {
	objectAsMap := map[string]any{}
	for _, mapR := range sliceRep {
		key := mapR["Key"].(string)
		val := mapR["Value"]
		if key == "_id" {
			key = "id"
		}

		if nestedData, ok := val.([]any); ok {
			nestedMaps := []map[string]any{}
			// cast all to map[string]any
			for _, nested := range nestedData {
				nestedMaps = append(nestedMaps, nested.(map[string]any))
			}
			objectAsMap[key] = db.getObjectAsMap(nestedMaps)
		} else {
			objectAsMap[key] = val
		}
	}
	return objectAsMap
}

// a field of "" and value of "" will return all records in the collection
func (db *MongoWrapper) GetRecordsByField(field string, value any) ([]map[string]any, error) {
	var results []any

	var filter bson.D
	if field == "" && value == "" {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{Key: field, Value: value}}
	}
	cursor, err := db.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &results)

	if err != nil {
		return nil, err
	}

	jsonRep, err := json.Marshal(results)

	if err != nil {
		return nil, err
	}

	var sliceOfsliceRep = [][]map[string]any{}

	err = json.Unmarshal(jsonRep, &sliceOfsliceRep)

	if err != nil {
		return nil, err
	}

	objectsMapList := []map[string]any{}
	for _, objectDesc := range sliceOfsliceRep {
		objectAsMap := db.getObjectAsMap(objectDesc)
		objectsMapList = append(objectsMapList, objectAsMap)
	}

	return objectsMapList, nil
}

func (db *MongoWrapper) GetIdByFieldAndValue(field string, value any) string {
	records, err := db.GetRecordsByField(field, value)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return ""
	}
	if len(records) == 0 {
		return ""
	}

	if len(records) > 1 {
		panic("MongoWrapper.GetIdByFieldAndValue: length of records returned should be 1")
	}

	return records[0]["id"].(string)
}

func (db *MongoWrapper) GetAllOfRecords() []map[string]any {
	records, _ := db.GetRecordsByField("", "")
	return records
}

func (db *MongoWrapper) Delete(id string) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	db.collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: objectId}})
}

func (db *MongoWrapper) Update(id string, data UpdateDesc) bool {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	result, err := db.collection.UpdateByID(context.Background(),
		objectId,
		bson.D{{Key: "$set", Value: bson.D{{Key: data.Field, Value: data.Value}}}})
	if err != nil {
		return false
	}

	return result.ModifiedCount == 1
}

func (db *MongoWrapper) DeleteDb() error {
	return db.client.Database(db.database_name).Drop(context.Background())
}

func (db *MongoWrapper) Commit() error {
	return nil
}

var MONGO_WRAPPER_MAP = map[string]*MongoWrapper{}

func MakeMongoWrapper(database string, collection string) (*MongoWrapper, error) {

	if collection == "" {
		panic("MakeMongoWrapper: collection cannot be empty")
	}

	if database == "" {
		panic("MakeMongoWrapper: database cannot be empty")
	}

	key := database + collection
	// implements singleton pattern
	if MONGO_WRAPPER_MAP[key] != nil {
		return MONGO_WRAPPER_MAP[key], nil
	}

	file_db, err := new(MongoWrapper).New(database, collection)

	if err != nil {
		panic("MakeMongoWrapper: " + err.Error())
	}

	MONGO_WRAPPER_MAP[key] = file_db
	return file_db, nil
}

func RemoveMongoSingleton(database, collection string, shouldDeleteDatabase ...bool) {
	if collection == "" {
		panic("RemoveMongoSingleton: collection cannot be empty")
	}

	if database == "" {
		panic("RemoveMongoSingleton: database cannot be empty")
	}

	key := database + collection

	MongoEng, exists := MONGO_WRAPPER_MAP[key]
	if exists {
		// this is just to make deleting a database dificult and intentional
		if len(shouldDeleteDatabase) == 1 {
			if shouldDeleteDatabase[0] && database == "test" {
				MongoEng.DeleteDb()
			}
		}

		if len(shouldDeleteDatabase) == 2 {
			if shouldDeleteDatabase[0] && shouldDeleteDatabase[1] {
				MongoEng.DeleteDb()
			}
		}

		delete(MONGO_WRAPPER_MAP, key)
	}
}
