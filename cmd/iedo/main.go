package main

import (
	"fmt"

	"github.com/blck-snwmn/iedo"
)

func main() {
	if err := iedo.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
