package main

import (
	"go_utils/utils"
	"os"
	"os/exec"
	"syscall"
)

func RunClash() bool {
	procSync.lock.Lock()
	defer procSync.lock.Unlock()
	if procSync.state {
		return false
	}
	go func() {
		c := exec.Command("bash", "-c", clashBinary+" -d "+clashConfigDir)
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
		err = c.Process.Signal(syscall.SIGTERM)
		if err != nil {
			utils.LogPrintError("Failed send SIGTERM")
			utils.LogPrintError(err)
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
