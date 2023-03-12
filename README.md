# go-random

[![CI](https://github.com/susisu/go-random/workflows/CI/badge.svg)](https://github.com/susisu/go-random/actions?query=workflow%3ACI)

go-random provides generic random number generator interfaces and functions.

## Usage

Use `go get` to install:

``` shell
go get github.com/susisu/go-random
```

go-random provides functions for both uint32 and uint64 generators.
The functions for each type are exported from separate pacakges, so use the appropriate one for your use case.

Here is an example using the uint64 version:

``` go
import (
	"math/rand"

	random "github.com/susisu/go-random/uint64"
)

func main() {
	// math/rand.Source64 implements the uint64 version of random.Generator
	g := rand.NewSource(42).(rand.Source64)
	// use go-random to generate random values of variaous numeric types
	v := random.Float64(g)
	fmt.Printf("%f\n", v)
}
```

## License

[MIT License](http://opensource.org/licenses/mit-license.php)

## Author

Susisu ([GitHub](https://github.com/susisu), [Twitter](https://twitter.com/susisu2413))
