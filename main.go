package main

import (
	"docker-server/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.NewCmdRoot().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Done!")
	os.Exit(0)
}
