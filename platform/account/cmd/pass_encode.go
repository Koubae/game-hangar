package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	str := "manu:123"
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	fmt.Println(encoded) // Output: bWFudToxMjM=
}
