package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	IrInterval int `json:"VentState, number"`
}

/*Config contains user preferences*/
type Config config

/*
ToJSON returns serialized date
*/
func (q *Config) ToJSON() (d []byte, err error) {
	//now := int(time.Now().Unix())

	b, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (q *Config) saveConfig() error {
	b, err := json.Marshal(q)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ioutil.WriteFile(configFileName, b, 0644)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (q *Config) setDefault() {
	q.IrInterval = 60 //secs
}

func (q *Config) loadConfig() error {

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Println("Cannot open config: ", err.Error())
		q.setDefault()
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(q); err != nil {
		log.Println("Error parsing config: ", err.Error())
		q.setDefault()
		return err
	}

	return nil
}

func notifyConfigChanged() error {
	wh.ReportWsEvent("configChanged", "1")
	return nil
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	defer notifyConfigChanged()
	d, errs := conf.ToJSON()
	if errs != nil {
		log.Println(errs.Error())
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//io.WriteString(w, string(d))

	state := r.FormValue("state")
	if len(state) == 0 {
		log.Println("config requested:")
		io.WriteString(w, string(d))
		return
	}
	_, errs = strconv.Atoi(state)
	if errs != nil {
		log.Println(errs.Error())
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
}
