package kvlist

import (
	"reflect"
	"testing"
)

func TestKeyValueList_Add(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key2", "value2"})

	if l.Count() != 2 {
		t.Error("expected two items")
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value1"})) {
		t.Error("invalid item 0")
	}

	if !reflect.DeepEqual(*l.list[1], (KeyValue{"key2", "value2"})) {
		t.Error("invalid item 1")
	}
}

func TestKeyValueList_Put(t *testing.T) {
	l := new(KeyValueList).Put(KeyValue{"key1", "value1"}).Put(KeyValue{"key1", "value2"})

	if l.Count() != 1 {
		t.Error("expected one item")
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value2"})) {
		t.Error("invalid item 0")
	}
}

func TestKeyValueList_DeleteKey_Tail(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key2", "value2"})

	if !l.DeleteKey("key2") {
		t.Error("key2 not deleted")
	}

	if l.Count() != 1 {
		t.Error("expected one item")
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value1"})) {
		t.Error("invalid item 0")
	}
}

func TestKeyValueList_DeleteKey_Head(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key2", "value2"})

	if !l.DeleteKey("key1") {
		t.Error("key2 not deleted")
	}

	if l.Count() != 1 {
		t.Error("expected one item")
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key2", "value2"})) {
		t.Error("invalid item 0")
	}
}

func TestKeyValueList_DeleteKeys_Found(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key1", "value2"})

	if !l.DeleteKeys("key1") {
		t.Error("key1 not deleted")
	}

	if l.Count() != 0 {
		t.Error("expected no items")
	}
}

func TestKeyValueList_DeleteKeys_NotFound(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key1", "value2"})

	if l.DeleteKeys("key2") {
		t.Error("key2 deleted")
	}

	if l.Count() != 2 {
		t.Error("expected two items")
	}
}

func TestKeyValueList_Get(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key2", "value2"})

	if !reflect.DeepEqual(l.Get(0), (KeyValue{"key1", "value1"})) {
		t.Error("invalid item 0")
	}

	if !reflect.DeepEqual(l.Get(1), (KeyValue{"key2", "value2"})) {
		t.Error("invalid item 1")
	}
}

func TestKeyValueList_GetKey_Found(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"})

	if key, ok := l.GetKey("key1"); ok {
		if !reflect.DeepEqual(key, (KeyValue{"key1", "value1"})) {
			t.Error("invalid item 0")
		}
	} else {
		t.Error("key1 not found")
	}
}

func TestKeyValueList_GetKey_NotFound(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"})

	if _, ok := l.GetKey("key2"); ok {
		t.Error("key2 found")
	}
}

func TestKeyValueList_GetKeys_Found(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key1", "value2"})
	kvs := l.GetKeys("key1")

	if !reflect.DeepEqual(kvs, []KeyValue{{"key1", "value1"}, {"key1", "value2"}}) {
		t.Error("not all items found")
	}
}

func TestKeyValueList_GetKeys_NotFound(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key1", "value2"})
	kvs := l.GetKeys("key2")

	if len(kvs) != 0 {
		t.Error("unexpected items found")
	}
}

func TestKeyValueList_Items(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key1", "value2"})

	if !reflect.DeepEqual(l.Items(), []KeyValue{{"key1", "value1"}, {"key1", "value2"}}) {
		t.Error("not all items found")
	}
}

func TestKeyValueList_String(t *testing.T) {
	l := new(KeyValueList).Add(KeyValue{"key1", "value1"}).Add(KeyValue{"key2", "value2"})

	if l.String() != "[key1=value1 key2=value2]" {
		t.Error("KeyValueList.String() failed")
	}
}