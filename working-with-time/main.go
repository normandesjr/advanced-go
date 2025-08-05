package main

import (
	"fmt"
	"time"
)

func main() {

	now := time.Now().UTC()
	fmt.Println(now.Format("2006-01-02T15:04:05.000Z"))
}
