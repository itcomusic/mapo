# Map Order

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![coverage-img]][coverage-url]

Is a custom implementation of a map in Go that maintains the order of keys as they are added
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
    
    b, _ := json.Marshal(m)
    fmt.Println(string(b)) // {"a":1,"b":2}
    
    var m2 mapo.Map
    _ = json.Unmarshal(b, &m2)
    fmt.Println(m2.Keys()) // [a, b]
}

```

## License

[MIT License](LICENSE)

[build-img]: https://github.com/itcomusic/mapo/workflows/build/badge.svg

[build-url]: https://github.com/itcomusic/mapo/actions

[pkg-img]: https://pkg.go.dev/badge/github.com/itcomusic/mapo.svg

[pkg-url]: https://pkg.go.dev/github.com/itcomusic/mapo

[coverage-img]: https://codecov.io/gh/itcomusic/mapo/branch/main/graph/badge.svg

[coverage-url]: https://codecov.io/gh/itcomusic/mapo