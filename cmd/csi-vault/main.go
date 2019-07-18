package main

import (
	"os"

	"kmodules.xyz/client-go/logs"
	"kubevault.dev/csi-driver/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
