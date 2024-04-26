package main

import (
	"encoding/json"
	"go_tools/cmd/g2/event"
	"go_utils/utils"
	"net/http"
)

func handleUserInput(w http.ResponseWriter, req *http.Request) {
	gasData := req.FormValue("gas_data")
	e := &event.PlayerEvent{}
	if err := json.Unmarshal([]byte(gasData), e); err != nil {
		utils.LogPrintError("Cannot parse user input")
		return
	}
	GlobalEventChannel <- e
	<-e.Finished
	w.Write([]byte(e.Output))
}
