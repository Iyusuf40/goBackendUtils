package tests

import (
	"testing"
	"time"

	"github.com/Iyusuf40/goBackendUtils/storage"
)

var RS storage.TempStore
var TEST_DB = 15

func beforeEachRSF() {
	RS = storage.MakeRedisWrapper(TEST_DB)
}

func afterEachRSF() {
	rs := RS.(*storage.RedisWrapper)
	rs.FlushDB()
}

func TestSetKeyToValAndGetValRS(t *testing.T) {
	beforeEachRSF()
	defer afterEachRSF()

	key := "key"
	RS.DelKey(key)
	got := RS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty, got:", got)
	}

	val := "value"

	ok := RS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestSetKeyToValAndGetVal: failed to set value")
	}
	got = RS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be " + val + " got " + got)
	}
}

func TestDelKeyRS(t *testing.T) {
	beforeEachRSF()
	defer afterEachRSF()

	key := "key2"
	val := "value"

	ok := RS.SetKeyToVal(key, val)
	if !ok {
		t.Fatal("TestDelKey: failed to set value")
	}
	got := RS.GetVal(key)
	// key exists
	if got != val {
		t.Fatal("TestDelKey: expected value to be " + val + " got " + got)
	}

	ok = RS.DelKey(key)
	if !ok {
		t.Fatal("TestDelKey: failed to delete value")
	}

	got = RS.GetVal(key)
	// key does not exists
	if got != "" {
		t.Fatal("TestDelKey: expected value to be empty")
	}
}

func TestSetKeyToValWithExpiryRS(t *testing.T) {
	beforeEachRSF()
	defer afterEachRSF()

	key := "key3"
	RS.DelKey(key)
	got := RS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"
	expiry := 1.5

	ok := RS.SetKeyToValWIthExpiry(key, val, expiry)
	if !ok {
		t.Fatal("TestSetKeyToValWIthExpiry: failed to set value")
	}

	got = RS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// value should still exists after half duration
	time.Sleep(time.Second * time.Duration(expiry/2))
	got = RS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// value should not exist after duration
	time.Sleep(time.Second * time.Duration(int(expiry+1)))
	got = RS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be empty, got " + got)
	}
}

func TestChangeKeyExpiryRS(t *testing.T) {
	beforeEachRSF()
	defer afterEachRSF()

	key := "key4"
	RS.DelKey(key)
	got := RS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValAndGetVal: expected value to be empty")
	}

	val := "value"
	expiry := 1.5

	ok := RS.SetKeyToValWIthExpiry(key, val, expiry)
	if !ok {
		t.Fatal("TestSetKeyToValWIthExpiry: failed to set value")
	}

	got = RS.GetVal(key)
	if got != val {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be " + val + " got " + got)
	}

	// shorten the expiry
	RS.ChangeKeyEpiry(key, expiry/2)

	// value should not exist after half duration, since
	// it has been shortened
	time.Sleep(time.Second * time.Duration(expiry/2))
	got = RS.GetVal(key)
	if got != "" {
		t.Fatal("TestSetKeyToValWIthExpiry: expected value to be empty, got " + got)
	}
}
