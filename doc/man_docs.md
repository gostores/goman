# Generating Man Pages For Your Own goman.Command

Generating man pages from a goman command is incredibly easy. An example is as follows:

```go
package main

import (
	"log"

	"github.com/govenue/goman"
	"github.com/govenue/goman/doc"
)

func main() {
	cmd := &goman.Command{
		Use:   "test",
		Short: "my test program",
	}
	header := &doc.GenManHeader{
		Title: "MINE",
		Section: "3",
	}
	err := doc.GenManTree(cmd, header, "/tmp")
	if err != nil {
		log.Fatal(err)
	}
}
```

That will get you a man page `/tmp/test.3`
