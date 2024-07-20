package mapo

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestMap_Set(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		mo := New()
		mo.Set("a", 34)
		mo.Set("b", []int{3, 4, 5})

		if got, want := mo.key, []string{"a", "b"}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}

		if got, want := mo.value, map[string]any{
			"a": 34,
			"b": []int{3, 4, 5},
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}
	})

	t.Run("copy key", func(t *testing.T) {
		t.Parallel()

		mo := New()
		mo.Set("a", 34)
		mo.Set("b", []int{3, 4, 5})
		mo.Set("a", 35)

		if got, want := mo.key, []string{"a", "b"}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}

		if got, want := mo.value, map[string]any{
			"a": 35,
			"b": []int{3, 4, 5},
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}
	})

	t.Run("nil map", func(t *testing.T) {
		t.Parallel()

		var mo Map
		mo.Set("a", 34)

		if got, want := mo.key, []string{"a"}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}

		if got, want := mo.value, map[string]any{
			"a": 34,
		}; !reflect.DeepEqual(got, want) {
			t.Errorf("%v not equal to expected %v", got, want)
		}
	})
}

func TestNewMapWithSize(t *testing.T) {
	t.Parallel()

	mo := NewWithSize(2)
	if got, want := cap(mo.key), 2; got != want {
		t.Errorf("cap: %v not equal to expected %v", got, want)
	}
}

func TestMap_Get(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		mo := New()
		mo.Set("a", 34)

		got, ok := mo.Get("a")
		if want := true; ok != want {
			t.Fatalf("%v not equal to expected %v", ok, want)
		}

		if want := 34; got != want {
			t.Errorf("%v not equal to expected %v", got, want)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mo := New()
		_, ok := mo.Get("a")
		if want := false; ok != want {
			t.Fatalf("%v not equal to expected %v", ok, want)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var mo Map
		_, ok := mo.Get("a")
		if want := false; ok != want {
			t.Fatalf("%v not equal to expected %v", ok, want)
		}
	})
}

func TestMap_Keys(t *testing.T) {
	t.Parallel()

	mo := New()
	mo.Set("a", 34)
	mo.Set("b", []int{3, 4, 5})

	if got, want := mo.Keys(), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Errorf("%v not equal to expected %v", got, want)
	}
}

func TestMap_Delete(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		mo := New()
		mo.Set("a", 34)
		mo.Delete("a")

		want := &Map{
			value: make(map[string]any),
			key:   []string{},
		}
		if !reflect.DeepEqual(mo, want) {
			t.Errorf("%+v not equal to expected %+v", mo, want)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mo := New()
		mo.Set("a", 34)
		mo.Delete("b")

		want := &Map{
			value: map[string]any{"a": 34},
			key:   []string{"a"},
		}
		if !reflect.DeepEqual(mo, want) {
			t.Errorf("%+v not equal to expected %+v", mo, want)
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	mo := New()
	mo.Set("a", 34)
	mo.Set("b", []int{3, 4, 5})

	got, err := json.Marshal(mo)
	if err != nil {
		t.Fatalf("%v", err)
	}

	want := `{"a":34,"b":[3,4,5]}`
	if !bytes.Equal(got, []byte(want)) {
		t.Errorf("%q not equal to expected %q", got, want)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	in := `{
		"a": 1,
		"b": "hello",
		"c": true,
		"d": false,
		"e": null,
		"f": 3.14,
		"g": [1, 2, 3],
		"h": {
				"i": 4,
				"j": "world"
			}
		}`

	want := Map{
		value: map[string]any{
			"a": json.Number("1"),
			"b": "hello",
			"c": true,
			"d": false,
			"e": nil,
			"f": json.Number("3.14"),
			"g": []any{json.Number("1"), json.Number("2"), json.Number("3")},
			"h": &Map{
				value: map[string]any{
					"i": json.Number("4"),
					"j": "world",
				},
				key: []string{"i", "j"},
			},
		},
		key: []string{"a", "b", "c", "d", "e", "f", "g", "h"},
	}

	var got Map
	if err := json.Unmarshal([]byte(in), &got); err != nil {
		t.Fatalf("%v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%+v not equal to expected %+v", got, want)
	}
}

func TestUnmarshalJSON_Invalid(t *testing.T) {
	t.Parallel()

	tt := []struct {
		in string
	}{
		{in: ""},
		{in: "[]"},
		{in: "["},
		{in: "{]"},
		{in: "{3"},
		{in: "{}3"},
		{in: "{3:"},
		{in: `{"key": }"`},
		{in: `{"key": 3, "b": [{`},
		{in: `{"key": 3, "b": [}`},
	}

	for _, tc := range tt {
		t.Run(tc.in, func(t *testing.T) {
			t.Parallel()

			mo := New()
			err := mo.UnmarshalJSON([]byte(tc.in))
			if err == nil {
				t.Fatalf("expecting error: %v", tc.in)
			}
		})
	}
}

func TestMarshalJSON_Invalid(t *testing.T) {
	t.Parallel()

	om := New()
	om.Set("m", math.NaN())
	if _, err := json.Marshal(om); err == nil {
		t.Fatalf("expecting error: %+v", om)
	}
}
