package gevent

import (
	"time"
)

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IStore interface {
		Get() *Event
		Put(*Event)
	}

	defaultStore struct {
		eventChan chan *Event //event buffer channel
		backCount int         //maximum backup buffer
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize the default store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewDefaultStore(backCount int) IStore {
	return &defaultStore{
		eventChan: make(chan *Event, backCount),
		backCount: backCount,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get event from store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultStore) Get() *Event {
	var event *Event

	if s.backCount == 1 {
		event = <-s.eventChan
	} else {
		select {
		case eventSource := <-s.eventChan:
			event = eventSource
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	return event
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set event to store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultStore) Put(event *Event) {
	s.eventChan <- event
}
