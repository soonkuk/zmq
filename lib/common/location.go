package common

import "errors"

type Location struct {
	latitude  float32
	longitude float32
}

func NewLocation(lat float32, longi float32) (*Location, error) {
	if (lat) < 0 {
		return nil, errors.New("wrong lattitude")
	}
	if (longi) < 0 {
		return nil, errors.New("wrong longitude")
	}

	location := &Location{
		latitude:  lat,
		longitude: longi,
	}

	return location, nil
}
