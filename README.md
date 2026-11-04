# mailmap

A Go library for parsing Git `.mailmap` files and resolving author identities. It uses only the standard library and supports the four mailmap entry forms for replacing names, email addresses, or both.

## Installation

```sh
go get github.com/git-pkgs/mailmap
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/git-pkgs/mailmap"
)

func main() {
	m, err := mailmap.Parse(strings.NewReader(`Jane Doe <jane@example.com> <jane@laptop.(none)>
`))
	if err != nil {
		log.Fatal(err)
	}
	name, email := m.Resolve("jane", "jane@laptop.(none)")
	fmt.Printf("%s <%s>\n", name, email)
	// Jane Doe <jane@example.com>
}
```

`Parse` accepts an `io.Reader`, including an open `.mailmap` file. `New()` creates an empty map; empty, zero-value, and nil maps return identities unchanged.

Matching ignores name and email case. Entries matching both name and email take priority over entries matching only email. Within either group, the first matching entry wins. An unmatched identity is returned unchanged.

Blank lines, comments starting with `#`, and malformed entries are skipped. Reader errors and lines exceeding the default `bufio.Scanner` token limit return an error along with the entries parsed so far.

## License

[MIT](LICENSE).
