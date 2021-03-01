package bucket

import (
	"fmt"
	"time"
)

func buildState(id ID, stream []Event) (state, error) {

	ret := state{}

	if len(stream) == 0 {
		return ret, fmt.Errorf("Empty stream")
	}

	for _, e := range stream {

		if id != e.StreamID() {
			return state{}, fmt.Errorf("ID Mismatch")
		}

		switch event := e.(type) {
		case *Opened:
			ret.id = event.id
			ret.name = event.name
			ret.desc = event.desc
			ret.v = event.v
			ret.last = event.at
		case *Updated:
			ret.name = event.name
			ret.desc = event.desc
			ret.v = event.v
			ret.last = event.at
		case *Closed:
			ret.closed = true
			ret.v = event.v
			ret.last = event.at
		default:
			return state{}, fmt.Errorf("Stream contains unkown events")
		}
	}

	return ret, nil
}

type state struct {
	id     ID
	name   Name
	desc   Description
	closed bool
	last   time.Time
	v      Version
}
