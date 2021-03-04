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
		t:    "bucket.Open",
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

func (s *store) InsertEvent(ctx context.Context, e events.Event) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.checkVersion(e)
}

func (s *store) checkVersion(e events.Event) error {
	const op errors.Op = "inmem.store.checkVersion"

	var expect events.EntityVersion

	daos, ok := s.data[e.EntityID()]
	if !ok {
		expect = 1
	} else {
		expect = (daos[len(daos)-1].v + 1)
	}

	if expect != e.EntityVersion() {
		return errors.New(op, errors.KindUnexpected, "version error")
	}

	return s.encodeToDao(e)
}

func (s *store) encodeToDao(e events.Event) error {
	const op errors.Op = "inmem.store.encodeToDao"

	var d dao

	d.encode(e)

	return s.insert(e.EntityID(), d)
}

func (s *store) insert(id events.EntityID, d dao) error {

	s.data[id] = append(s.data[id], d)

	return nil

}

func (s *store) GetStream(ctx context.Context, id events.EntityID) ([]events.Event, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.findStream(id)
}

func (s *store) findStream(id events.EntityID) ([]events.Event, error) {
	const op errors.Op = "inmem.store.findStream"

	daos, ok := s.data[id]
	if !ok {
		return nil, errors.New(op, errors.KindNotFound, "Stream not found")
	}

	return decodeToEvents(id, daos)
}

func decodeToEvents(id events.EntityID, daos []dao) ([]events.Event, error) {
	const op errors.Op = "inmem.decodeToEvents"

	ret := make([]events.Event, len(daos)+1)

	for _, dao := range daos {

		e, err := dao.decode(id)
		if err != nil {
			return nil, errors.New(op, errors.KindUnexpected, "decoding error", err)
		}
		ret = append(ret, e)
	}

	return ret, nil
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
		return &bucket.Opened{eb, data}, nil

	case "bucket.Updated":
		data := d.data.(bucket.BucketData)
		return &bucket.Updated{eb, data}, nil

	case "bucket.Closed":
		return &bucket.Closed{eb}, nil

	default:
		return nil, errors.New(op, errors.KindUnexpected, "Unkown event type: "+d.t)
	}
}
