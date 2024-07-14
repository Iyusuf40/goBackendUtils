package tests

import (
	"os"
	"testing"
	"time"

	"github.com/Iyusuf40/goBackendUtils/storage"
)

var temp_test_db_path = "temp_store_test_db.json"
var TS storage.TempStore

func beforeEachTSF() {
	TS = storage.MakeTempStoreFileDbImpl(temp_test_db_path, "T")
}

func afterEachTSF() {
	os.Remove(temp_test_db_path)
	storage.RemoveDbSingleton(temp_test_db_path, "T")
}

func TestSetKeyToValAndGetVal(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key"
	TS.DelKey(key)
	got := TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"

	ok := TS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestSetKeyToValAndGetVal: failed to set value")
	}
	got = TS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be " + val + " got " + got)
	}
}

func TestDelKey(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key2"
	val := "value"

	ok := TS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestDelKey: failed to set value")
	}
	got := TS.GetVal(key)
	// key exists
	if got != val {
		t.Fatal("TestDelKey: expected value to be " + val + " got " + got)
	}

	ok = TS.DelKey(key)
	if !ok {
		t.Fatal("TestDelKey: failed to delete value")
	}

	got = TS.GetVal(key)
	// key does not exists
	if got != "" {
		t.Fatal("TestDelKey: expected value to be empty")
	}
}

func TestSetKeyToValWithExpiry(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key3"
	TS.DelKey(key)
	got := TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"
	expiry := 1.5

	ok := TS.SetKeyToValWIthExpiry(key, val, expiry)
	if !ok {
		t.Fatal("TestSetKeyToValWIthExpiry: failed to set value")
	}

	got = TS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// value should still exists after half duration
	time.Sleep(time.Second * time.Duration(expiry/2))
	got = TS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// value should not exist after duration
	time.Sleep(time.Second * time.Duration(int(expiry)))
	got = TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be empty, got " + got)
	}
}

func TestChangeKeyExpiry(t *testing.T) {
	beforeEachTSF()
	defer afterEachTSF()

	key := "key4"
	TS.DelKey(key)
	got := TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"
	expiry := 1.5

	ok := TS.SetKeyToValWIthExpiry(key, val, expiry)
	if !ok {
		t.Fatal("TestSetKeyToValWIthExpiry: failed to set value")
	}

	got = TS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// shorten the expiry
	TS.ChangeKeyEpiry(key, expiry/2)

	// value should not exist after half duration, since
	// it has been shortened
	time.Sleep(time.Second * time.Duration(expiry/2))
	got = TS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be empty, got " + got)
	}
}
