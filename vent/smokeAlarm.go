package vent

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

var firstTimeKitchen = true
var firstTimeSauna = true
var alarmTimeSauna int32
var alarmTimeKitchen int32

var saunaAlarmStatus = false
var kitchenAlarmStatus = false
var r = raspi.NewRaspiAdaptor("raspi")

/*
SmokeAlarmKitchen smoke/gas sensor in the kitchen
*/
var SmokeAlarmKitchen = gpio.NewButtonDriver(r, "alm1", "11") //17

/*
SmokeAlarmSauna smoke/gas sensor in the sauna
*/
var SmokeAlarmSauna = gpio.NewButtonDriver(r, "alm2", "13") //27

func init() {
	log.Println("Initialise smoke alarms")

	gobot.On(SmokeAlarmKitchen.Event("release"), func(data interface{}) {
		now := int32(time.Now().Unix())
		//log.Println(now - alarmTimeKitchen)
		if (now-alarmTimeKitchen) < 2 || kitchenAlarmStatus == true {
			return
		}
		alarmTimeKitchen = now
		kitchenAlarmStatus = true

		log.Println("Smoke alarm in the kitchen: On")

		/*
			err := Sms("Test", "Smoke/Gas detected in the Kitchen", Recipients)
			if err != nil {
				log.Println(err.Error())
			}
		*/

		ReportAlert("Something is really not OK", "Smoke/Gas detected in the Kitchen")

	})

	gobot.On(SmokeAlarmKitchen.Event("push"), func(data interface{}) {
		if firstTimeKitchen {
			firstTimeKitchen = false
			return
		}

		now := int32(time.Now().Unix())
		if (now-alarmTimeKitchen) < 2 || kitchenAlarmStatus == false {
			return
		}
		alarmTimeKitchen = now
		kitchenAlarmStatus = false

		log.Println("Smoke alarm in the kitchen: Off")

		/*
			err := Sms("Test", "Smoke/Gas in the Kitchen - all good", Recipients)
			if err != nil {
				log.Println(err.Error())
			}
		*/

		ReportAlert("Now it's better", "Smoke/Gas in the Kitchen - all good")

	})

	gobot.On(SmokeAlarmSauna.Event("release"), func(data interface{}) {
		now := int32(time.Now().Unix())
		//log.Println(now - alarmTimeKitchen)
		if now-alarmTimeSauna < 2 || saunaAlarmStatus == true {
			return
		}

		alarmTimeSauna = now
		saunaAlarmStatus = true

		log.Println("Smoke alarm in the Sauna: On")
		/*
			err := Sms("Test", "Smoke/Gas detected in the Sauna", Recipients)
			if err != nil {
				log.Println(err.Error())
			}
		*/

		ReportAlert("Something is really not OK", "Smoke/Gas detected in the Sauna")

	})

	gobot.On(SmokeAlarmSauna.Event("push"), func(data interface{}) {
		if firstTimeSauna {
			firstTimeSauna = false
			return
		}

		now := int32(time.Now().Unix())
		//log.Println(now - alarmTimeKitchen)
		if now-alarmTimeSauna < 2 || saunaAlarmStatus == false {
			return
		}

		alarmTimeSauna = now
		saunaAlarmStatus = false

		log.Println("Smoke alarm in the Sauna: Off")
		/*
			err := Sms("Test", "Smoke/Gas in Sauna - all good", Recipients)
			if err != nil {
				log.Println(err.Error())
			}
		*/
		ReportAlert("Now it's better", "Smoke/Gas in Sauna - all good")

	})
}
