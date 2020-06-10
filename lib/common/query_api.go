package common

import (
	"errors"
	"math/rand"
	"time"

	"github.com/vmihailenco/msgpack"
)

type QueryResponser struct {
	data      Query
	QueryChan chan Query
	status    ResponserStatus
}

func NewQueryResponser(status ResponserStatus) *QueryResponser {
	d := NewQuery()
	return &QueryResponser{data: d, status: status}
}

func (qr *QueryResponser) Run(c chan Query) {
	// qr.QueryChan = c
	var temperatureRange []float32
	if qr.status == CorrectResponser {
		temperatureRange = []float32{35.0, 37.5}
	} else {
		temperatureRange = []float32{37.5, 42.0}
	}

	defer close(c)
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		qr.data.Temperature = temperatureRange[0] + rand.Float32()*(temperatureRange[1]-temperatureRange[0])
		qr.data.Hrv = rand.Intn(33) + 49
		latitude := 35.13 + rand.Float32()*2.63
		longitude := 35.13 + rand.Float32()*2.63
		qr.data.Location, _ = NewLocation(latitude, longitude)
		qr.data.TimeStamp = time.Now().UTC()
		c <- qr.data
	}
}

type Query struct {
	Location    location  `json:"location"`
	Hrv         int       `json:"hrv"`
	Ecg         int16     `json:"ecg"`
	Temperature float32   `json:"temperature"`
	TimeStamp   time.Time `json:"timeStamp"`
}

func NewQuery() Query {
	return Query{}
}

func (q *Query) ToJson() ([]byte, error) {
	return msgpack.Marshal(q)
}

type location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func NewLocation(lat float32, longi float32) (location, error) {
	var loc location
	if (lat) < 0 {
		return loc, errors.New("wrong lattitude")
	}
	if (longi) < 0 {
		return loc, errors.New("wrong longitude")
	}

	loc = location{
		Latitude:  lat,
		Longitude: longi,
	}

	return loc, nil
}

type ResponserStatus string

const (
	CorrectResponser ResponserStatus = "CorrectResponser"
	FailResponser    ResponserStatus = "FailResponser"
)
