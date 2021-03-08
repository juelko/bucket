package inmem

import (
	"context"
	"sync"

	"github.com/juelko/bucket/bucket"

	"github.com/juelko/bucket/pkg/errors"
	"github.com/juelko/bucket/pkg/events"
)

func NewBucketStore() bucket.Store {
	return &store{
		mtx:  sync.RWMutex{},
		data: map[events.EntityID][]dao{},
	}
}

func NewTestBucketStore() bucket.Store {
	open := dao{
		t:    "bucket.Opened",
		v:    1,
		data: bucket.BucketData{Title: "OpenTitle", Description: "Open Description"},
	}
	updated := dao{
		t:    "bucket.Updated",
		v:    2,
		data: bucket.BucketData{Title: "UpdatedTitle", Description: "Updated Description"},
	}
	closed := dao{
		t:    "bucket.Closed",
		v:    3,
		data: nil,
	}

	return &store{
		mtx: sync.RWMutex{},
		data: map[events.EntityID][]dao{
			"OpenID":    {open},
			"UpdatedID": {open, updated},
			"ClosedID":  {open, updated, closed},
		},
	}
}

type store struct {
	mtx  sync.RWMutex
	data map[events.EntityID][]dao
}

func (s *store) OpenStream(ctx context.Context, o *bucket.Opened) error {
	const op errors.Op = "inmem.store.OpenStream"
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.exists(o.EntityID()) {
		return errors.New(op, errors.KindAllreadyExists, "Allready exists")
	}

	return s.insert(o)

}

func (s *store) InsertEvent(ctx context.Context, e events.Event) error {
	const op errors.Op = "inmem.store.InsertEvent"
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if !s.exists(e.EntityID()) {
		return errors.New(op, errors.KindNotFound, "Stream not found")
	}

	if s.nextVersion(e.EntityID()) != e.EntityVersion() {
		return errors.New(op, errors.KindUnexpected, "version error")
	}

	return s.insert(e)
}

func (s *store) insert(e events.Event) error {
	var d dao

	d.encode(e)

	s.data[e.EntityID()] = append(s.data[e.EntityID()], d)

	return nil
}

func (s *store) GetStream(ctx context.Context, id events.EntityID) ([]events.Event, error) {
	const op errors.Op = "inmem.store.GetStream"
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	daos, ok := s.data[id]
	if !ok {
		return []events.Event{}, errors.New(op, errors.KindNotFound, "Stream not found")
	}

	return decodeToEvents(id, daos)
}

func decodeToEvents(id events.EntityID, daos []dao) ([]events.Event, error) {
	const op errors.Op = "inmem.decodeToEvents"

	ret := make([]events.Event, len(daos))

	for i, dao := range daos {

		e, err := dao.decode(id)
		if err != nil {
			return []events.Event{}, errors.New(op, errors.KindUnexpected, "decoding error", err)
		}
		ret[i] = e
	}

	return ret, nil
}

func (s *store) exists(id events.EntityID) bool {

	_, ok := s.data[id]

	return ok
}

func (s *store) nextVersion(id events.EntityID) events.EntityVersion {
	return events.EntityVersion(len(s.data[id]) + 1)
}

// data access object
type dao struct {
	t    string // value from events.Event.Type()
	v    events.EntityVersion
	data interface{}
}

func (d *dao) encode(e events.Event) {
	d.t = e.Type()
	d.v = e.EntityVersion()
	d.data = e.Data()
}

// decodes dao and given id to returned event,
// returns error and nil if type string dao.ts has unknown value
func (d dao) decode(id events.EntityID) (events.Event, error) {
	const op errors.Op = "inmem.dao.build"

	eb := events.Base{
		ID: id,
		V:  d.v,
	}

	switch d.t {
	case "bucket.Opened":
		data := d.data.(bucket.BucketData)
		return &bucket.Opened{Base: eb, BucketData: data}, nil

	case "bucket.Updated":
		data := d.data.(bucket.BucketData)
		return &bucket.Updated{Base: eb, BucketData: data}, nil

	case "bucket.Closed":
		return &bucket.Closed{Base: eb}, nil

	default:
		return nil, errors.New(op, errors.KindUnexpected, "Unkown event type: "+d.t)
	}
}
