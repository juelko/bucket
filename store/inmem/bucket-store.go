package inmem

import (
	"context"
	"sync"
	"time"

	"github.com/juelko/bucket/bucket"
	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
	"github.com/juelko/bucket/pkg/request"
)

// data access object
type dao struct {
	t    string
	at   time.Time
	v    events.Version
	rid  request.ID
	data []byte
}

func (d dao) base(id events.StreamID) events.Base {
	return events.Base{
		ID:  id,
		RID: d.rid,
		At:  d.at,
		V:   d.v,
	}
}

func NewBucketStore() bucket.Store {
	return &store{
		mtx:  sync.RWMutex{},
		data: map[events.StreamID][]dao{},
	}
}

type store struct {
	mtx  sync.RWMutex
	data map[events.StreamID][]dao
}

func (s *store) InsertEvent(ctx context.Context, e events.Event) error {
	const op errors.Op = "inmem.InsertEvent"

	return s.insert(e)
}

func (s *store) insert(e events.Event) error {
	const op errors.Op = "inmem.insert"

	s.mtx.Lock()
	defer s.mtx.Unlock()

	var expect events.Version

	daos, ok := s.data[e.StreamID()]
	if !ok {
		expect = 1
	} else {
		expect = (daos[len(daos)-1].v + 1)
	}

	if expect != e.Version() {
		return errors.New(op, errors.KindUnexpected, "Version error")
	}

	dao := dao{
		t:    e.Type(),
		at:   e.Occured(),
		v:    e.Version(),
		rid:  e.RequestID(),
		data: e.Payload(),
	}

	s.data[e.StreamID()] = append(daos, dao)

	return nil
}

func (s *store) GetStream(ctx context.Context, id events.StreamID) ([]events.Event, error) {
	const op errors.Op = "inmem.GetStream"
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	daos, ok := s.data[id]
	if !ok {
		return nil, errors.New(op, errors.KindNotFound, "Stream not found")
	}

	return buildEvents(id, daos)
}

func buildEvents(id events.StreamID, daos []dao) ([]events.Event, error) {
	const op errors.Op = "inmem.buildEvents"

	ret := []events.Event{}

	for _, dao := range daos {

		e, err := bucket.BuildEvent(dao.base(id), dao.t, dao.data)
		if err != nil {
			return nil, err
		}
		ret = append(ret, e)
	}

	return ret, nil
}
