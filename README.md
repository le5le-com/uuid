# le5le-com/uuid

[![Go Reference](https://pkg.go.dev/badge/github.com/le5le-com/uuid.svg)](https://pkg.go.dev/github.com/le5le-com/uuid)

A UUID package for Go, support for converting mongodb objectID to uuid.

It currently only supports UUID v7.

## Installation

```bash
go get github.com/le5le-com/uuid
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/le5le-com/uuid"
)

func main() {
	uuidv7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	fmt.Printf("UUIDv7: %s, time=%v\n", uuidv7, uuidv7.TimeFromV7())
  // UUIDv7: 0189dd43-c284-7f4f-806e-e7d238e9babb，time=2023-08-10 10:25:52.772 +0800 CST

  s := "0189dd43-c284-7f4f-806e-e7d238e9babb"
	uuid, err := Parse(s)
	if err != nil {
	  panic(err)
	}

  objectId := "63ede45a8d0137fc1b631091"
	uuidv7, err = UUIDV7FromObjectID(objectId)
	if err != nil {
		panic(err)
	}

  if uuidv7.ObjectIDHex() != objectId {
		fmt.Printf("Convert uuidv7 to objectId error: uuid.ObjectIDHex=%s, objectId=%s", uuid.ObjectIDHex(), objectId)
	}
}

```

## License

[MIT](./LICENSE.md)
