package vent

import (
	"io"
	"log"
	"net/http"
)

/*CurrentState contains current app data*/
var CurrentState HData

func init() {
	CurrentState = HData{}
}

/*
CurrentStateHandler - Reports current state to sockets
*/
func CurrentStateHandler(w http.ResponseWriter, r *http.Request) {

	d, errs := CurrentState.ToJSON()
	if errs != nil {
		log.Println(errs.Error())
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(d))
}
