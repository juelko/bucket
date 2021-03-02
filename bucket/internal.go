package bucket

import (
	"fmt"
	"time"

	"github.com/juelko/bucket/pkg/events"
)

func buildState(id events.StreamID, stream []events.Event) (state, error) {

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
			ret.id = event.StreamID()
			ret.title = event.Title
			ret.desc = event.Desc
			ret.v = event.Version()
			ret.last = event.Occured()
		case *Updated:
			ret.title = event.Title
			ret.desc = event.Desc
			ret.v = event.Version()
			ret.last = event.Occured()
		case *Closed:
			ret.closed = true
			ret.v = event.Version()
			ret.last = event.Occured()
		default:
			return state{}, fmt.Errorf("Stream contains unkown events")
		}
	}

	return ret, nil
}

type state struct {
	id     events.StreamID
	title  Title
	desc   Description
	closed bool
	last   time.Time
	v      events.Version
}
