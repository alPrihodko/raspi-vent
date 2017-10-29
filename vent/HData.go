package vent

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

/*
LIMIT max amount of history data values
*/
const LIMIT = 20

/*
HData is set of home data values, can be parsed from json
*/
type HData struct {
	KitchenSensorState bool `json:"KitchenSensorState, boolean"`
	SaunaSensorState   bool `json:"SaunaSensorState, boolean"`

	KitchenVentState bool  `json:"FlowerBadState, boolean"`
	KitchenVentTimer int64 `json:"FlowerBadTimer, Number"`

	Timestamp int `json:"Timestamp, int"`
	Index     int `json:"index, int"`
}

/*
HistoryData is storage for recent states
*/
type HistoryData []*HData

/*
Len of HistoryData
*/
func (q HistoryData) Len() int { return len(q) }

/*
Push HomeData to HistoryData
*/
func (q *HistoryData) Push(x interface{}) {
	n := len(*q)
	old := *q
	if n >= LIMIT {
		item := old[0]
		item.Index = -1 // for safety
		*q = old[1:n]
		item = nil
	}
	n = len(*q)
	item := x.(*HData)
	item.Index = n
	*q = append(*q, item)
}

/*
Pop HomeData from HistoryData
*/
func (q *HistoryData) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

/*
Last - returns last record
*/
func (q *HistoryData) Last() HData {
	old := *q
	n := len(old)
	if n == 0 {
		return HData{}
	}
	item := old[n-1]
	return *item
}

/*
ToJSON returns serialized hash
*/
func (q *HData) ToJSON() (d []byte, err error) {
	//now := int(time.Now().Unix())

	b, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/*
ToJSON returns serialized hash
*/
func (q *HistoryData) ToJSON(from int) (d []byte, err error) {
	old := *q
	sl := HistoryData{}
	//now := int(time.Now().Unix())
	log.Println("loading history starting from: ", from)
	if from > 0 {
		var interval = 1

		for i := 0; i < old.Len(); i = i + interval {
			index := i
			if i >= old.Len() {
				index = old.Len() - 1
			}

			item := old[index]
			if item.Timestamp > from {
				sl.Push(item)
			}
		}
	} else {
		sl = old
	}

	b, err := json.Marshal(sl)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/*
ToGOB64 encodes to string
*/
func (q *HistoryData) ToGOB64() (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&q)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

/*
FromGOB64 decodes from string
*/
func (q *HistoryData) FromGOB64(str string) error {
	//q := &HistoryData{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
		return err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(q)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return nil
}

/*
SerializeToFile writes slice to file
*/
func (q *HistoryData) SerializeToFile(name string) error {
	f, err := os.Create(name)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	w := bufio.NewWriter(f)
	str, err := q.ToGOB64()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	x, err := w.WriteString(str)
	log.Println("bytes written: ", x)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	w.Flush()
	return nil
}

/*
RestoreFromFile writes slice to file
*/
func (q *HistoryData) RestoreFromFile(name string) error {
	dat, err := ioutil.ReadFile(name)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = q.FromGOB64(bytes.NewBuffer(dat).String())

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

/*
HistoryDataHandler - sends it over http
*/
func (q *HistoryData) HistoryDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	from, errx := strconv.Atoi(r.URL.Query().Get("from"))
	if errx != nil {
		log.Println(errx.Error())
		from = 0
	}

	d, errs := q.ToJSON(from)
	if errs != nil {
		log.Println(errs.Error())
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(d))
}
