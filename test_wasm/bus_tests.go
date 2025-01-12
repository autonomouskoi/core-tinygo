package main

import (
	"bytes"
	"fmt"

	"github.com/autonomouskoi/akcore"
	bus "github.com/autonomouskoi/core-tinygo"
)

func main() {}

//go:export start
func Start() int32 {
	TestKVSetGetDelete()
	TestKVList()
	return 0
}

//go:export recv
func Recv() int32 {
	return 0
}

func SendMessage(testName, error string) {
	msg := &bus.BusMessage{
		Topic:   "HOST",
		Message: []byte(testName),
	}
	if error != "" {
		msg.Error = &bus.Error{Detail: &error}
	}
	bus.Send(msg)
}

func TestKVSetGetDelete() {
	testName := "KVSetGetDelete"
	k := []byte("test-key")
	v := []byte("test-value")
	// key doesn't exist yet
	if _, err := bus.KVGet(k); err == nil {
		SendMessage(testName, "expected error on non-existent key")
		return
	}
	// set it
	if err := bus.KVSet(k, v); err != nil {
		SendMessage(testName, "error setting key: "+err.Error())
		return
	}
	// get it
	got, err := bus.KVGet(k)
	if err != nil {
		SendMessage(testName, "error getting key: "+err.Error())
		return
	}
	// make sure it matches
	if !bytes.Equal(got, v) {
		SendMessage(testName, fmt.Sprintf("got %q, want %q", string(got), string(v)))
		return
	}
	// delete it
	if err := bus.KVDelete(k); err != nil {
		SendMessage(testName, "error deleting key: "+err.Error())
		return
	}
	// make sure it's gone
	if _, err := bus.KVGet(k); err == nil {
		SendMessage(testName, "expected error on deleted key")
		return
	} else if err != akcore.ErrNotFound {
		SendMessage(testName, "expected ErrNotFound, got "+err.Error())
		return
	}
	SendMessage(testName, "")
}

func TestKVList() {
	testName := "KVList"
	v := []byte("test-value")
	// create some test values
	for i := 0; i < 4; i++ {
		k := []byte(fmt.Sprint("a", i))
		if err := bus.KVSet(k, v); err != nil {
			SendMessage(testName, "setting key "+string(k))
			return
		}
		k = []byte(fmt.Sprint("b", i))
		if err := bus.KVSet(k, v); err != nil {
			SendMessage(testName, "setting key "+string(k))
			return
		}
	}
	values, err := bus.KVList([]byte("b"), 2, 1)
	if err != nil {
		SendMessage(testName, "listing values: "+err.Error())
		return
	}
	if values.GetTotalMatches() != 4 {
		SendMessage(testName, fmt.Sprint("wanted 4 matches, got ", values.GetTotalMatches()))
		return
	}
	if len(values.GetKeys()) != 2 {
		SendMessage(testName, fmt.Sprint("wanted 2 keys, got ", len(values.GetKeys())))
		return
	}
	if !bytes.Equal(values.GetKeys()[0], []byte("b1")) {
		SendMessage(testName, fmt.Sprintf(`first key want "b1", got %s`, string(values.GetKeys()[0])))
		return
	}
	if !bytes.Equal(values.GetKeys()[1], []byte("b2")) {
		SendMessage(testName, fmt.Sprintf(`first key want "b2", got %s`, string(values.GetKeys()[1])))
		return
	}
	SendMessage(testName, "")
}
