package bucket

import (
	"time"
)

func newState(stream []Event) state {
	ret := state{}

	for _, e := range stream {
		switch event := e.(type) {
		case Opened:
			ret.id = event.id
			ret.name = event.Name
			ret.description = event.Description
			ret.version = event.version
			ret.updated = event.occured
		case Updated:
			ret.name = event.Name
			ret.description = event.Description
			ret.version = event.version
			ret.updated = event.occured
		case Closed:
			ret.closed = true
			ret.version = event.version
			ret.updated = event.occured
		}
	}

	return ret
}

type state struct {
	id          BucketID
	name        string
	description string
	version     uint
	closed      bool
	updated     time.Time
}
