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
	IPublishSubscribe interface {
		Subscribe(ISubscriberHandler, ...int) IPublishSubscribe
		Publish(IEventSource) IPublishSubscribe
	}

	IEventHub interface {
		GetChannel(channelName string, args ...ITransport) IChannel
		Broadcast(data interface{}) IEventHub
	}

	eventHub struct {
		channels map[string]IChannel //channel collection
		mu       sync.Mutex          //synchronous
	}
)

var (
	hub     IEventHub
	hubOnce sync.Once
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get event hub instance
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetEventHub() IEventHub {
	hubOnce.Do(func() {
		hub = &eventHub{
			channels: make(map[string]IChannel, 0),
		}
	})

	return hub
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get channels
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) GetChannel(channelName string, args ...ITransport) IChannel {
	s.mu.Lock()
	defer s.mu.Unlock()

	channel, isOk := s.channels[channelName]
	if !isOk {
		var transport ITransport
		if len(args) > 0 {
			transport = args[0]
		} else {
			transport = NewDefaultTransport(1024)
		}

		channel = NewChannel(channelName, transport)
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
