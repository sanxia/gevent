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
	IEventCenter interface {
		IEvent
		GetEventHub() IEventHub
	}

	eventCenter struct {
		hub IEventHub
	}
)

var (
	center     IEventCenter
	centerOnce sync.Once
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取IEventCenter实例
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func GetEventCenterInstance() IEventCenter {
	centerOnce.Do(func() {
		center = &eventCenter{
			hub: GetEventHubInstance(),
		}
	})

	return center
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取IEventHub
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventCenter) GetEventHub() IEventHub {
	return s.hub
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * publish event
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventCenter) Publish(eventSource IEventSource) IEvent {
	if eventSource != nil {
		if channelName := eventSource.GetChannelName(); len(channelName) > 0 {
			if channel := s.hub.GetChannel(channelName); channel != nil {
				if eventName := eventSource.GetEventName(); len(eventName) > 0 {
					channel.Publish(eventName, eventSource.GetData())
				}
			}
		}
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * subscribe event
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *eventCenter) Subscribe(subscriberHandler ISubscriberHandler, args ...int) IEvent {
	priority := 0
	repeat := 0
	argsCount := len(args)

	if argsCount > 0 {
		priority = args[0]
	}

	if argsCount > 1 {
		repeat = args[1]
	}

	if subscriberHandler != nil {
		if channelName := subscriberHandler.GetChannelName(); len(channelName) > 0 {
			if channel := s.hub.GetChannel(channelName); channel != nil {
				if eventName := subscriberHandler.GetEventName(); len(eventName) > 0 {
					channel.Subscribe(eventName, func(e *Event) {
						subscriberHandler.Handler(e)
					}, priority, repeat)
				}
			}
		}
	}

	return s
}
