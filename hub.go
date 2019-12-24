package gevent

import (
	"sync"
)

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IEventHub interface {
		GetChannel(channelName string, args ...IStore) IChannel
		Broadcast(data interface{}) IEventHub
	}

	eventHub struct {
		channels  map[string]IChannel //channel collection
		backCount int                 //maximum backup buffer
		mu        sync.Mutex          //synchronous
	}
)

var (
	hub     IEventHub
	hubOnce sync.Once
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get event hub instance
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetEventHub(backCount int) IEventHub {
	hubOnce.Do(func() {
		hub = newEventHub(backCount)
	})

	return hub
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize event hub
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func newEventHub(backCount int) IEventHub {
	eventDispatcher := &eventHub{
		channels:  make(map[string]IChannel, 0),
		backCount: backCount,
	}

	return eventDispatcher
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get channels
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) GetChannel(channelName string, args ...IStore) IChannel {
	s.mu.Lock()
	defer s.mu.Unlock()

	channel, isOk := s.channels[channelName]
	if !isOk {
		var store IStore
		if len(args) > 0 {
			store = args[0]
		} else {
			store = NewDefaultStore(s.backCount)
		}

		channel = NewChannel(channelName, store)
		s.channels[channelName] = channel
	}

	return channel
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * broadcast events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) Broadcast(data interface{}) IEventHub {
	for _, channel := range s.channels {
		channel.Broadcast(data)
	}

	return s
}
