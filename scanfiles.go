package main

import (
	"fmt"
	"wattpad_challenge/scanutils"
)

func main() {
	err, filename := scanutils.RunScanFile()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Text Scan file successfully created:", filename)
	}
}
