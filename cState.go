package main

import (
	"log"
	"raspi-vent/vent"
	"raspi-vent/wsHandler"
	"time"
)

func appStateChanged() {
	log.Println("app state change triggered")

	//irr.CurrentState.GardenName = ir01.Relay.Name()
	//log.Println("registering mode: " + ir01.GetMode())
	//irr.CurrentState.GardenMode, irr.CurrentState.GardenTimer = ir01.GetMode()
	//irr.CurrentState.GardenState = ir01.GetState()

	vent.CurrentState.Timestamp = int(time.Now().Unix())

	/*update UI*/
	d, errs := vent.CurrentState.ToJSON()
	if errs != nil {
		vent.ReportAlert(errs.Error(), "Cannot report current state")
		return
	}

	err := wsHandler.ReportData(d)
	if err != nil {
		vent.ReportAlert(err.Error(), "Cannot report relay state")
	}

	x := vent.CurrentState

	historyData.Push(&x)
}
