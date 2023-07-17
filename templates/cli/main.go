package main

import (
	"flag"
	"go_utils/utils"
)

func main() {
	var verboseFlag = flag.Bool("v", false, "verbose mode")
	flag.Parse()

	if *verboseFlag {
		utils.SetLogLevelDebug()
	} else {
		utils.SetLogLevelInfo()
	}
	utils.LogPrintInfo("Log level: ", utils.LogLevel)

	// otherArgs := flag.Args()

}
