package main

import (
	"fmt"
	"os"

	"github.com/moomaideng/bogbank/services/ledger/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
