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
		RegisterChannel(channels ...IChannel) IEventHub
		GetChannel(channelName string) IChannel
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
func GetEventHubInstance() IEventHub {
	hubOnce.Do(func() {
		hub = &eventHub{
			channels: make(map[string]IChannel, 0),
		}
	})

	return hub
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * register channel
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) RegisterChannel(channels ...IChannel) IEventHub {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, channel := range channels {
		if currentChannel, isOk := channel.(*Channel); isOk {
			if _, isExists := s.channels[currentChannel.Name]; !isExists {
				s.channels[currentChannel.Name] = channel
			}
		}
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get channels
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventHub) GetChannel(channelName string) IChannel {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channel, isOk := s.channels[channelName]; isOk {
		return channel
	}

	return nil
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
