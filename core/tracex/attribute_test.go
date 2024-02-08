package tracex

import (
	"go.opentelemetry.io/otel/attribute"
	"testing"
)

func TestCreateAttribute(t *testing.T) {

	// 1. attribute string
	attr := attribute.String("key1", "value1")
	if attr.Key != "key" {
		t.Errorf("Expected key 'key', got '%s'", attr.Key)
	}

	if attr.Value.AsString() != "value" {
		t.Errorf("Expected value 'value', got '%s'", attr.Value.AsString())
	}

	// 2. attribute KeyValue
	attr2 := attribute.KeyValue{TraceName, attribute.StringValue("hello-world")}
	if attr2.Key != TraceName {
		t.Errorf("Expected key 2 'key', got '%s'", attr2.Key)
	}

	if attr2.Value.AsString() != "value" {
		t.Errorf("Expected value 2 'value', got '%s'", attr2.Value.AsString())
	}

	// 3.
	// ....
}

// 测试 attribute.Key.Bool() 方法是否正确创建属性
func TestKeyBool(t *testing.T) {
	key := attribute.Key("test_key")
	value := true
	kv := key.Bool(value)

	if kv.Key != key || kv.Value.Type() != attribute.BOOL || kv.Value.AsBool() != value {
		t.Errorf("Key.Bool() = %v, want %v", kv.Value, value)
	}
}

// 测试 attribute.Key.Int64() 方法是否正确创建属性
func TestKeyInt64(t *testing.T) {
	key := attribute.Key("test_key")
	value := int64(42)
	kv := key.Int64(value)

	if kv.Key != key || kv.Value.Type() != attribute.INT64 || kv.Value.AsInt64() != value {
		t.Errorf("Key.Int64() = %v, want %v", kv.Value, value)
	}
}

// 测试 attribute插入到attribute.Set后是否正确
func TestAttributeSetInsert(t *testing.T) {
	key := attribute.Key("test_key")
	value := "test_value"
	kv := key.String(value)

	set := attribute.NewSet(kv)
	gotValue, ok := set.Value(key)

	if !ok || gotValue.AsString() != value {
		t.Errorf("Attribute set did not contain the expected value, got %v, want %v", gotValue, value)
	}
}

func TestAttributeSetIter(t *testing.T) {
	key1 := attribute.Key("test_key1")
	value1 := "test_value1"
	kv1 := key1.String(value1)

	key2 := attribute.Key("test_key2")
	value2 := int64(42)
	kv2 := key2.Int64(value2)

	set := attribute.NewSet(kv1, kv2)
	iter := set.Iter()

	found := make(map[string]attribute.KeyValue)
	for iter.Next() {
		kv := iter.Attribute()
		found[string(kv.Key)] = kv
	}

	if len(found) != 2 || found[string(key1)].Value.AsString() != value1 ||
		found[string(key2)].Value.AsInt64() != value2 {
		t.Errorf("Iterator did not return the expected attributes")
	}
}
