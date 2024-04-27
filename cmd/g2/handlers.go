package main

import (
	"encoding/json"
	"go_tools/cmd/g2/event"
	"go_utils/utils"
	"net/http"
)

type ReturnData struct {
}

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

func handleGetData(w http.ResponseWriter, req *http.Request) {
	userID := req.FormValue("user_id")
	cityName := req.FormValue("city_name")
	if userID == "" {
		utils.LogPrintError("Invalid user id")
		w.Write([]byte("Invalid user id"))
		return
	}
	if cityName == "" {
		utils.LogPrintError("Invalid city name")
		w.Write([]byte("Invalid city name"))
		return
	}
	if _, ok := GlobalPlayerList[userID]; !ok {
		utils.LogPrintError("User not found")
		w.Write([]byte("User not found"))
		return
	}
	p := GlobalPlayerList[userID]
	if _, ok := p.CityList[cityName]; !ok {
		utils.LogPrintError("City not found")
		w.Write([]byte("City not found"))
		return
	}

}
