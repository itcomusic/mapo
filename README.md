# mapo

Map Order (Mapo) is a custom implementation of a map in Go that maintains the order of keys as they are added.
## Features

- **Order Preservation**: Keys are stored in the order they are added
- **Custom JSON Marshalling**: Implements custom JSON marshalling and unmarshalling to ensure the order of keys is preserved when converting to and from JSON
- **Key Management**: Provides methods to add, delete, and retrieve keys in their maintained order

## Installation

Go version 1.21+

```bash
go get github.com/itcomusic/mapo
```

## Usage
```go
import (
    "fmt"

    "github.com/itcomusic/mapo"
)

func main() {
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
}

```