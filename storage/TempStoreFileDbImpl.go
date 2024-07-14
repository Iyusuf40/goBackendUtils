package storage

import (
	"reflect"
	"time"
)

type TempStoreFileDbImpl struct {
	db       *FileDb
	MapStore map[string]any
	TimerMap map[string]any
	Init     bool
}

func (TS *TempStoreFileDbImpl) New(dbFile, recordsName string) TempStore {
	if recordsName == "" {
		recordsName = reflect.TypeOf(*TS).Name()
	}

	fileDb, _ := MakeFileDb(dbFile, recordsName)
	TS.db = fileDb
	TS.reload()
	return TS
}

func (TS *TempStoreFileDbImpl) reload() {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		TS.MapStore = map[string]any{}
		TS.TimerMap = map[string]any{}
		TS.Init = true
		TS.db.Save(*TS)
		TS.commit()
	} else {
		TS.runVacuum()
	}
}

func (TS *TempStoreFileDbImpl) GetVal(key string) string {
	mapStore := TS.getMapStore()

	if mapStore == nil {
		return ""
	}

	if TS.isExpired(key) {
		TS.deleteKeyTimerAndValue(key)
		return ""
	}

	val, _ := mapStore[key].(string)
	return val
}

func (TS *TempStoreFileDbImpl) SetKeyToVal(key string, value string) bool {
	mapStore := TS.getMapStore()
	if mapStore == nil {
		return false
	}
	mapStore[key] = value
	TS.commit()
	return true
}

// expiry is in seconds
func (TS *TempStoreFileDbImpl) setKeyToExpiry(key string, expiry float64) bool {
	timerStore := TS.getTimerMap()
	expiryTime := time.Now().Add(time.Second * time.Duration(expiry)).Unix()
	timerStore[key] = float64(expiryTime)
	TS.commit()
	return true
}

// sets key to val in MapSTore
// sets key to expiry in TimerMap
func (TS *TempStoreFileDbImpl) SetKeyToValWIthExpiry(key string, value string, expiry float64) bool {
	TS.setKeyToExpiry(key, expiry)
	TS.SetKeyToVal(key, value)
	return true
}

func (TS *TempStoreFileDbImpl) ChangeKeyEpiry(key string, newExpiry float64) bool {
	TS.setKeyToExpiry(key, newExpiry)
	return true
}

func (TS *TempStoreFileDbImpl) DelKey(key string) bool {
	mapStore := TS.getMapStore()

	if TS.keyExistsInTimerMap(key) {
		TS.deleteKeyTimerAndValueHelper(key)
	} else {
		delete(mapStore, key)
	}
	TS.commit()
	return true
}

func (TS *TempStoreFileDbImpl) DeleteDb() {
	TS.db.DeleteDb()
}

func (TS *TempStoreFileDbImpl) commit() {
	TS.db.Commit()
}

// runs throug all keys in TimerMap and checks ones that have expired
// removes the expired keys and their values from TimerMap and MapStore
func (TS *TempStoreFileDbImpl) runVacuum() {
	timerMap := TS.getTimerMap()
	changed := false
	for key := range timerMap {
		if TS.isExpired(key) {
			TS.deleteKeyTimerAndValueHelper(key)
			changed = true
		}
	}

	if changed {
		TS.commit()
	}
}

func (TS *TempStoreFileDbImpl) isExpired(key string) bool {
	if TS.keyExistsInTimerMap(key) {
		unixNow := time.Now().Unix()
		timeToExpire := TS.getTimerMap()[key].(float64)
		if timeToExpire <= float64(unixNow) {
			return true
		}

	}
	return false
}

func (TS *TempStoreFileDbImpl) deleteKeyTimerAndValue(key string) {
	TS.deleteKeyTimerAndValueHelper(key)
	TS.commit()
}

func (TS *TempStoreFileDbImpl) deleteKeyTimerAndValueHelper(key string) {
	timerMap := TS.getTimerMap()
	delete(timerMap, key)
	mapStore := TS.getMapStore()
	delete(mapStore, key)
}

func (TS *TempStoreFileDbImpl) getMapStore() map[string]any {
	mapStoreList, _ := TS.db.GetRecordsByField("Init", true)
	if mapStoreList == nil {
		return nil
	}

	if len(mapStoreList) > 1 {
		panic(`TempStoreFileDbImpl: getMapStore: there should only be one 
			instance in existence`)
	}
	mapStore := mapStoreList[0]["MapStore"].(map[string]any)
	return mapStore
}

func (TS *TempStoreFileDbImpl) getTimerMap() map[string]any {
	mapStoreList, _ := TS.db.GetRecordsByField("Init", true)
	if mapStoreList == nil {
		return nil
	}

	if len(mapStoreList) > 1 {
		panic(`TempStoreFileDbImpl: getTimerMap: there should only be one 
			instance in existence`)
	}
	timerMap := mapStoreList[0]["TimerMap"].(map[string]any)
	return timerMap
}

func (TS *TempStoreFileDbImpl) keyExistsInTimerMap(key string) bool {
	timerMap := TS.getTimerMap()
	_, exists := timerMap[key]
	return exists
}

func MakeTempStoreFileDbImpl(db_path, recordsName string) TempStore {
	return new(TempStoreFileDbImpl).New(db_path, recordsName)
}
