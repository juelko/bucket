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
	ts     string // value from events.Event.Type()
	at     time.Time
	v      events.Version
	rid    request.ID
	title  bucket.Title
	desc   bucket.Description
	closed bool // flag for closed streams
}

func (d *dao) encode(e events.Event) error {
	const op errors.Op = "inmem.dao.parse"

	d.ts = e.Type()
	d.at = e.Occured()
	d.v = e.Version()
	d.rid = e.RequestID()
	switch t := e.(type) {
	case *bucket.Opened:
		d.title = t.Title
		d.desc = t.Desc
		return nil
	case *bucket.Updated:
		d.title = t.Title
		d.desc = t.Desc
		return nil
	case *bucket.Closed:
		d.closed = true
		return nil
	default:
		return errors.New(op, errors.KindUnexpected, "Unkown event type: "+d.ts)
	}
}

// decodes dao and given id to returned event,
// returns error and nil if type string dao.ts has unknown value
func (d dao) decode(id events.StreamID) (events.Event, error) {
	const op errors.Op = "inmem.dao.build"

	eb := events.Base{
		ID:  id,
		RID: d.rid,
		At:  d.at,
		V:   d.v,
	}

	switch d.ts {
	case "bucket.Opened":
		return &bucket.Opened{Base: eb, Title: d.title, Desc: d.desc}, nil

	case "bucket.Updated":
		return &bucket.Updated{Base: eb, Title: d.title, Desc: d.desc}, nil

	case "bucket.Closed":
		return &bucket.Closed{Base: eb}, nil

	default:
		return nil, errors.New(op, errors.KindUnexpected, "Unkown event type: "+d.ts)
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
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.checkVersion(e)
}

func (s *store) checkVersion(e events.Event) error {
	const op errors.Op = "inmem.store.checkVersion"

	var expect events.Version

	daos, ok := s.data[e.StreamID()]
	if !ok {
		expect = 1
	} else {
		expect = (daos[len(daos)-1].v + 1)
	}

	if expect != e.Version() {
		return errors.New(op, errors.KindUnexpected, "version error")
	}

	return s.encodeToDao(e)
}

func (s *store) encodeToDao(e events.Event) error {
	const op errors.Op = "inmem.store.encodeToDao"

	var d dao

	if err := d.encode(e); err != nil {
		return errors.New(op, errors.KindUnexpected, "encoding error")
	}

	return s.insert(e.StreamID(), d)
}

func (s *store) insert(id events.StreamID, d dao) error {

	s.data[id] = append(s.data[id], d)

	return nil

}

func (s *store) GetStream(ctx context.Context, id events.StreamID) ([]events.Event, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.findStream(id)
}

func (s *store) findStream(id events.StreamID) ([]events.Event, error) {
	const op errors.Op = "inmem.store.findStream"

	daos, ok := s.data[id]
	if !ok {
		return nil, errors.New(op, errors.KindNotFound, "Stream not found")
	}

	return decodeToEvents(id, daos)
}

func decodeToEvents(id events.StreamID, daos []dao) ([]events.Event, error) {
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
