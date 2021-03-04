package bucket

import (
	"fmt"

	"github.com/juelko/bucket/pkg/events"
)

func buildState(id events.EntityID, stream []events.Event) (state, error) {

	ret := state{}

	if len(stream) == 0 {
		return ret, fmt.Errorf("Empty stream")
	}

	for _, e := range stream {

		if id != e.EntityID() {
			return state{}, fmt.Errorf("ID Mismatch")
		}

		switch event := e.(type) {
		case *Opened:
			ret.id = event.EntityID()
			ret.title = event.Title
			ret.desc = event.Description
			ret.v = event.EntityVersion()
		case *Updated:
			ret.title = event.Title
			ret.desc = event.Description
			ret.v = event.EntityVersion()
		case *Closed:
			ret.closed = true
			ret.v = event.EntityVersion()
		default:
			return state{}, fmt.Errorf("Stream contains unkown events")
		}
	}

	return ret, nil
}

type state struct {
	id     events.EntityID
	title  Title
	desc   Description
	closed bool
	v      events.EntityVersion
}
