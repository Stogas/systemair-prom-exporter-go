package main

import (
	"fmt"
	"os"
)

func logModbusError(context string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", context, err)
}
