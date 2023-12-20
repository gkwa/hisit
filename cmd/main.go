package main

import (
	"os"

	"github.com/taylormonacelli/hisit"
)

func main() {
	code := hisit.Execute()
	os.Exit(code)
}
