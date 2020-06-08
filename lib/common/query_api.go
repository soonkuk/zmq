package common

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"
)

type QueryResponser struct {
	data      Query
	QueryChan chan Query
}

func NewQueryResponser() *QueryResponser {
	d := NewQuery()
	return &QueryResponser{data: d}
}

func (qr *QueryResponser) Run(c chan Query) {
	// qr.QueryChan = c
	defer close(c)
	for {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		qr.data.Temperature = 35.0 + rand.Float32()*7.0
		qr.data.Hrv = rand.Intn(33) + 49
		latitude := 35.13 + rand.Float32()*2.63
		longitude := 35.13 + rand.Float32()*2.63
		qr.data.Location, _ = NewLocation(latitude, longitude)
		c <- qr.data
	}
}

type Query struct {
	Location    location `json:"location"`
	Hrv         int      `json:"hrv"`
	Ecg         int16    `json:"ecg"`
	Temperature float32  `json:"temperature"`
}

func NewQuery() Query {
	return Query{}
}

func (q *Query) ToJson() ([]byte, error) {
	return json.Marshal(q)
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
