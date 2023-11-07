package main

import (
	"flag"
	"fmt"
	"go_utils/utils"
	"go_utils/utils/minikv"
	"os"
	"path"
	"slices"
)

var (
	kv       *minikv.KV
	scanPath string
)

func RecursiveListFiles(startPath string) ([]string, error) {
	fileList := []string{}

	pendingPath := []string{startPath}
	for {
		if len(pendingPath) == 0 {
			break
		}
		currentPath := pendingPath[0]
		utils.LogPrintDebug("Start scanning: ", currentPath)
		items, err := os.ReadDir(currentPath)
		if err != nil {
			utils.LogPrintError("Failed to list dir: ", currentPath, err)
			return []string{}, err
		}
		utils.LogPrintDebug2("Found items: ", items)
		for _, v := range items {
			if v.IsDir() {
				newPath := path.Join(currentPath, v.Name())
				utils.LogPrintDebug("Adding new path: ", newPath)
				pendingPath = append(pendingPath, newPath)
				continue
			}
			utils.LogPrintDebug("Adding new file: ", path.Join(currentPath, v.Name()))
			fileList = append(fileList, path.Join(currentPath, v.Name()))
		}
		utils.LogPrintDebug2("Pending list: ", pendingPath)
		pendingPath = slices.Delete(pendingPath, 0, 1)
		utils.LogPrintDebug2("Pending list: ", pendingPath)
	}
	return fileList, nil
}

func main() {
	var verboseFlag = flag.Int("v", 0, "debug (max 4)")
	var scanPathRaw = flag.String("d", ".", "path to scan")
	flag.Parse()

	utils.SetLogLevelByVerboseFlag(*verboseFlag)
	scanPath = *scanPathRaw

	kv = minikv.MustNewKV("logger", 0)
	kv.MustLoad()

	files, err := RecursiveListFiles(scanPath)
	utils.ErrExit(err)
	utils.LogPrintInfo(fmt.Sprintf("Found %d files under %s", len(files), scanPath))
	for _, v := range files {
		utils.LogPrintInfo("File: ", v)
	}
}
