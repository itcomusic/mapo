package mapo_test

import (
	"encoding/json"
	"fmt"

	"github.com/itcomusic/mapo"
)

func ExampleMap() {
	m := mapo.New()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	fmt.Println(m.Keys()) // [a, b, c]
	value, ok := m.Get("a")
	fmt.Println(value, ok) // 1, true
	m.Delete("c")

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b)) // {"a":1,"b":2}

	var m2 mapo.Map
	if err := json.Unmarshal(b, &m2); err != nil {
		panic(err)
	}
	fmt.Println(m2.Keys()) // [a, b]

	// Output:
	// [a b c]
	// 1 true
	// {"a":1,"b":2}
	// [a b]
}
