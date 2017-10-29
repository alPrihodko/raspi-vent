package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"raspi-vent/irRelay"
	"raspi-vent/vent"
	"raspi-vent/wsHandler"
	"syscall"
	"time"

	"github.com/hybridgroup/gobot"

	"golang.org/x/net/websocket"
)

const configFileName = "/etc/raspi-vent.conf"

/*
HISTORYDATASERIAL file which contains history data for my home
*/
const HISTORYDATASERIAL = "goRaspiVentData.b64"

/*
INTERVAL  Check sensors status with interval
*/
var INTERVAL int

//var err error

var conf Config

//var conns socketConns

//relays
var ir01 irRelay.Ir

//var ir02 irRelay.Ir
//var ir03 irRelay.Ir

var wh wsHandler.WsHandler

//var currentState vent.HData
var historyData vent.HistoryData

func main() {

	err := conf.loadConfig()
	if err != nil {
		log.Println("Likely use default configuration")
	}

	gbot := gobot.NewGobot()

	ir01 = irRelay.New("garden", "33", &wh, appStateChanged)
	http.HandleFunc("/control/"+ir01.Relay.Name(), ir01.RelayHandler)
	//ir02 = irRelay.New("flowerbad", "35", &wh, appStateChanged)
	//http.HandleFunc("/control/"+ir02.Relay.Name(), ir02.RelayHandler)
	//ir03 = irRelay.New("grapes", "31", &wh, appStateChanged)
	//http.HandleFunc("/control/"+ir03.Relay.Name(), ir03.RelayHandler)

	//ir03 = irRelay.New("grapes", "29", &wh, appStateChanged)
	//http.HandleFunc("/control/"+ir03.Relay.Name(), ir03.RelayHandler)

	flag.IntVar(&INTERVAL, "timeout", 60, "Timeout?")
	flag.Parse()

	log.Println("Timeout interval to track sensors: ", INTERVAL)
	historyData.RestoreFromFile(HISTORYDATASERIAL)
	vent.CurrentState = historyData.Last()
	http.Handle("/relays", websocket.Handler(relHandler))

	http.Handle("/", http.FileServer(http.Dir("ui")))
	http.HandleFunc("/control/currentState", vent.CurrentStateHandler)
	http.HandleFunc("/control/hdata", historyData.HistoryDataHandler)
	http.HandleFunc("/control/config", configHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGABRT)

	/*
		work := func() {
			//defer home.Stop()
			gobot.Every(time.Duration(INTERVAL)*time.Second, func() {
				log.Println("gobot heartbeat")
				//TODO: some reporting work here...
				//reportSensors(sensors)

			})

		}

		go func() {
			<-c
			log.Println("Save history data...")
			historyData.SerializeToFile(HISTORYDATASERIAL)
			irRelay.Stop()
			os.Exit(1)
		}()
		robot := gobot.NewRobot("blinkBot",
			[]gobot.Connection{vent.GetRelayAdaptor()},
			[]gobot.Device{vent.SmokeAlarmSauna, vent.SmokeAlarmKitchen},
			work,
		)
	*/

	//gbot.AddRobot(robot)

	go gbot.Start()

	stop := scheduleBackup(backupHistoryData, time.Duration(INTERVAL*60)*time.Second, &historyData, HISTORYDATASERIAL)

	err = http.ListenAndServe(":1236", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	if stop != nil {
		stop <- true
	}

}

func scheduleBackup(what func(*vent.HistoryData, string), delay time.Duration,
	q *vent.HistoryData, l string) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what(q, l)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func backupHistoryData(q *vent.HistoryData, local string) {
	historyData.SerializeToFile(local)

	if _, err := DB.UploadFile(local, "/backup/raspi-vent.b64", true, ""); err != nil {
		log.Printf("Error uploading %s: %s\n", local, err)
	} else {
		log.Printf("File %s successfully uploaded\n", local)
	}
}
