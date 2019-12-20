package gevent

import (
	"sync"
	"time"
)

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	eventHub struct {
		channels map[string]*Channel //channel collection
		maxCount int                 //maximum backup buffer
		mu       sync.Mutex          //synchronous
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize event hub
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewEventHub(maxCount int) *eventHub {
	dispatcher := &eventHub{
		channels: make(map[string]*Channel, 0),
		maxCount: maxCount,
	}

	return dispatcher
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get channels
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) GetChannel(channelName string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	channel, isOk := s.channels[channelName]
	if !isOk {
		channel = &Channel{
			Name:        channelName,
			eventChan:   make(chan *Event, s.maxCount),
			subscribers: make(map[string]SubscriberList, 0),
		}

		s.channels[channelName] = channel
	}

	if !channel.isRunning {
		channel.isRunning = true
		go channel.eventLoop()
	}

	return channel
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * broadcast events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) Broadcast(eventName string, data interface{}) *eventHub {
	for _, channel := range s.channels {
		go func(channel *Channel) {
			event := &Event{
				Name:         eventName,
				Data:         data,
				IsBroadcast:  true,
				CreationDate: time.Now().UnixNano(),
			}
			channel.dispatchEvent(event)
		}(channel)
	}

	return s
}
