package main

import (
	"go_utils/utils"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func RunClash() bool {
	procSync.lock.Lock()
	defer procSync.lock.Unlock()
	if procSync.state {
		return false
	}
	go func() {
		var c *exec.Cmd
		if runtime.GOOS == "windows" {
			c = exec.Command(`.\`+clashBinary+".exe", "-d", `.\`+clashConfigDir)
		} else if runtime.GOOS == "linux" {
			c = exec.Command("bash", "-c", clashBinary+" -d "+clashConfigDir)
		} else {
			utils.LogPrintInfo("Not supported!!!!")
			return
		}
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Start()
		if err != nil {
			utils.LogPrintError(err)
		}
		utils.LogPrintInfo("stared process")
		procSync.lock.Lock()
		procSync.state = true
		procSync.lock.Unlock()
		<-procChan
		utils.LogPrintWarning("Killing process")
		if runtime.GOOS == "windows" {
			err := c.Process.Kill()
			if err != nil {
				utils.LogPrintError("Failed kill proc in windows")
				utils.LogPrintError(err)
			}
		} else if runtime.GOOS == "linux" {
			err = c.Process.Signal(syscall.SIGTERM)
			if err != nil {
				utils.LogPrintError("Failed send SIGTERM")
				utils.LogPrintError(err)
			}
		}
		err = c.Wait()
		if err != nil {
			utils.LogPrintError("Failed wait, should be normal")
			utils.LogPrintError(err)
		}
		procSync.lock.Lock()
		procSync.state = false
		procSync.lock.Unlock()
	}()
	return true
}

func StopClash() bool {
	procSync.lock.Lock()
	defer procSync.lock.Unlock()
	if !procSync.state {
		return false
	}
	procChan <- 0
	return true
}
