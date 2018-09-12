package kvlist

import (
	"bytes"
	"reflect"
	"testing"
)

func TestKeyValue_Read(t *testing.T) {
	kv := KeyValue{"key1", `value 日本語 "quote"`}
	b := new(bytes.Buffer)

	if _, err := b.ReadFrom(&kv); err != nil {
		t.Fatal(err)
	}

	if b.String() != `key1="value \u65e5\u672c\u8a9e \"quote\""` {
		t.Fatal("values not matched")
	}
}

func TestKeyValue_Write(t *testing.T) {
	kv := KeyValue{}

	if _, err := kv.Write(bytes.NewBufferString(`   key1="value \u65e5\u672c\u8a9e \"quote\""  `).Bytes()); err != nil {
		t.Fatal(err)
	}

	if kv.Key != "key1" {
		t.Fatal("invalid key:", kv.Key)
	}

	if kv.Value != `value 日本語 "quote"` {
		t.Fatal("invalid value:", kv.Value)
	}
}

func TestKeyValueList_Add(t *testing.T) {
	l := New().Add("key1", "value1").Add("key2", "value2")

	if l.Count() != 2 {
		t.Fatalf("expected two items, got %d", l.Count())
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value1"})) {
		t.Fatal("invalid item 0")
	}

	if !reflect.DeepEqual(*l.list[1], (KeyValue{"key2", "value2"})) {
		t.Fatal("invalid item 1")
	}
}

func TestKeyValueList_Put(t *testing.T) {
	l := New().Put("key1", "value1").Put("key1", "value2")

	if l.Count() != 1 {
		t.Fatalf("expected one item, got %d", l.Count())
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value2"})) {
		t.Fatal("invalid item 0")
	}
}

func TestKeyValueList_DeleteKey_Tail(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value2"})

	if !l.DeleteKey("key2") {
		t.Fatal("key2 not deleted")
	}

	if l.Count() != 1 {
		t.Fatalf("expected one item, got %d", l.Count())
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", "value1"})) {
		t.Fatal("invalid item 0")
	}
}

func TestKeyValueList_DeleteKey_Head(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value2"})

	if !l.DeleteKey("key1") {
		t.Fatal("key2 not deleted")
	}

	if l.Count() != 1 {
		t.Fatalf("expected one item, got %d", l.Count())
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key2", "value2"})) {
		t.Fatal("invalid item 0")
	}
}

func TestKeyValueList_DeleteKeys_Found(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key1", "value2"})

	if !l.DeleteKeys("key1") {
		t.Fatal("key1 not deleted")
	}

	if l.Count() != 0 {
		t.Fatalf("expected no items, got %d", l.Count())
	}
}

func TestKeyValueList_DeleteKeys_NotFound(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key1", "value2"})

	if l.DeleteKeys("key2") {
		t.Fatal("key2 deleted")
	}

	if l.Count() != 2 {
		t.Fatalf("expected two items, got %d", l.Count())
	}
}

func TestKeyValueList_Clear(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value2"})

	l.Clear()

	if l.Count() != 0 {
		t.Fatalf("expected no items, got %d", l.Count())
	}
}

func TestKeyValueList_Get(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value2"})

	if !reflect.DeepEqual(l.Get(0), (KeyValue{"key1", "value1"})) {
		t.Fatal("invalid item 0")
	}

	if !reflect.DeepEqual(l.Get(1), (KeyValue{"key2", "value2"})) {
		t.Fatal("invalid item 1")
	}
}

func TestKeyValueList_GetKey_Found(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"})

	if key, ok := l.GetKey("key1"); ok {
		if !reflect.DeepEqual(key, ("value1")) {
			t.Fatal("invalid item 0")
		}
	} else {
		t.Fatal("key1 not found")
	}
}

func TestKeyValueList_GetKey_NotFound(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"})

	if _, ok := l.GetKey("key2"); ok {
		t.Fatal("key2 found")
	}
}

func TestKeyValueList_ScanKey_Found(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"})
	var v string

	if ok := l.ScanKey(&v, "key1"); ok {
		if v != "value1" {
			t.Fatal("invalid item 0")
		}
	} else {
		t.Fatal("key1 not found")
	}
}

func TestKeyValueList_ScanKey_NotFound(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"})
	var v string

	if ok := l.ScanKey(&v, "key2"); ok {
		t.Fatal("key2 found")
	}
}

func TestKeyValueList_GetKeys_Found(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key1", "value2"})
	kvs := l.GetKeys("key1")

	if !reflect.DeepEqual(kvs, []KeyValue{{"key1", "value1"}, {"key1", "value2"}}) {
		t.Fatal("not all items found")
	}
}

func TestKeyValueList_GetKeys_NotFound(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key1", "value2"})
	kvs := l.GetKeys("key2")

	if len(kvs) != 0 {
		t.Fatal("unexpected items found")
	}
}

func TestKeyValueList_Items(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key1", "value2"})

	if !reflect.DeepEqual(l.Items(), []KeyValue{{"key1", "value1"}, {"key1", "value2"}}) {
		t.Fatal("not all items found")
	}
}

func TestKeyValueList_String(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value2"})

	if l.String() != "[key1=value1 key2=value2]" {
		t.Fatal("KeyValueList.String() failed")
	}
}

func TestKeyValueList_Read(t *testing.T) {
	l := New().AddKeyValue(KeyValue{"key1", "value1"}).AddKeyValue(KeyValue{"key2", "value 日本語 \"quote\""})
	b := new(bytes.Buffer)

	if _, err := b.ReadFrom(l); err != nil {
		t.Fatal(err)
	}

	if b.String() != `key1="value1" key2="value \u65e5\u672c\u8a9e \"quote\""` {
		t.Fatal("values not matched")
	}
}

func TestKeyValueList_Write(t *testing.T) {
	s := `   key1="value \u65e5\u672c\u8a9e \"quote\""  key2="value 2"`
	l, err := NewFromString(s)

	if err != nil {
		t.Fatal(err)
	}

	if l.Count() != 2 {
		t.Fatalf("expected two items, got %d", l.Count())
	}

	if !reflect.DeepEqual(*l.list[0], (KeyValue{"key1", `value 日本語 "quote"`})) {
		t.Fatalf("invalid item 0")
	}

	if !reflect.DeepEqual(*l.list[1], (KeyValue{"key2", `value 2`})) {
		t.Fatal("invalid item 1")
	}
}
