package config

import "testing"

func TestReader(t *testing.T) {
	bb := []byte(`{"foo": "bar", "baz": {"bar": "cat"}}`)

	type data struct {
		path  []string
		value string
	}
	var datas []data
	datas = append(datas, data{
		path:  []string{"foo"},
		value: "bar",
	})
	datas = append(datas, data{
		path:  []string{"baz", "bar"},
		value: "cat",
	})

	r := NewReader()

	c, err := r.Merge(&ChangeSet{Data: bb}, &ChangeSet{})
	if err != nil {
		t.Fatal(err)
	}

	values, err := r.Values(c)
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range datas {
		if v := values.Get(test.path...).String(""); v != test.value {
			t.Fatalf("Expected %s got %s for path %v", test.value, v, test.path)
		}
	}
}
