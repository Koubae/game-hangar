package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	str := "admin:admin"
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	fmt.Println(encoded) // Output: bWFudToxMjM=
}
