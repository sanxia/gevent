package gevent

import (
	"sort"
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
	IChannel interface {
		Broadcast(data interface{}) IChannel
		Publish(eventName string, data interface{}) IChannel
		Subscribe(eventName string, eventHandler EventHandler, args ...int) IChannel
		Unsubscribe(eventNames ...string) IChannel
	}

	Channel struct {
		Name        string     //channel name
		transport   ITransport //event transport
		subscribers sync.Map   //subscriber map Collection
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * channel initialization
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewChannel(channelName string, transport ITransport) IChannel {
	channel := &Channel{
		Name:        channelName,
		transport:   transport,
		subscribers: sync.Map{},
	}

	go channel.eventLoop()

	return channel
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * broadcast events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Broadcast(data interface{}) IChannel {
	event := &Event{
		ChannelName:  s.Name,
		Name:         "__broadcast__",
		Data:         data,
		IsBroadcast:  true,
		CreationDate: time.Now().UnixNano(),
	}

	go s.dispatchEvent(event)

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * publishing events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Publish(eventName string, data interface{}) IChannel {
	event := &Event{
		ChannelName:  s.Name,
		Name:         eventName,
		Data:         data,
		CreationDate: time.Now().UnixNano(),
	}

	go s.dispatchEvent(event)

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * subscribe to events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Subscribe(eventName string, eventHandler EventHandler, args ...int) IChannel {
	subscribers := make(SubscriberList, 0)

	if datas, isOk := s.subscribers.Load(eventName); isOk {
		if datas, isOk := datas.(SubscriberList); isOk {
			subscribers = datas
		}
	}

	//优先级
	priority, repeat := 0, 0
	if len(args) > 0 {
		priority = args[0]
	}

	if len(args) > 1 {
		repeat = args[1]
	}

	subscribers = append(subscribers, &Subscriber{
		Priority:     priority,
		Repeat:       repeat,
		Counts:       sync.Map{},
		handler:      eventHandler,
		creationDate: time.Now().UnixNano(),
	})

	//优先级从高到低降序
	sort.Sort(sort.Reverse(subscribers))

	s.subscribers.Store(eventName, subscribers)

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * unsubscribe events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Unsubscribe(eventNames ...string) IChannel {
	if len(eventNames) == 0 {
		s.subscribers.Range(func(eventName, eventValue interface{}) bool {
			s.subscribers.Delete(eventName)
			return true
		})
	} else {
		for eventName := range eventNames {
			s.subscribers.Range(func(key, value interface{}) bool {
				if eventName == key {
					s.subscribers.Delete(eventName)
					return false
				}
				return true
			})
		}
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * dispatch event
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) dispatchEvent(event *Event) {
	if event.IsBroadcast {
		s.transport.Store(event)
	} else {
		if _, isOk := s.subscribers.Load(event.Name); isOk {
			s.transport.Store(event)
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event loop
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) eventLoop() {
	for {
		if event := s.transport.Load(); event != nil {
			s.subscribers.Range(func(eventName, eventValue interface{}) bool {
				if subscribers, isOk := eventValue.(SubscriberList); isOk {
					if eventName == event.Name {
						go s.eventHandler(event, subscribers)
					} else if event.IsBroadcast {
						go s.eventHandler(event, subscribers)
					}
				}

				return true
			})
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event handling
 * subscriber handler is sync notify
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) eventHandler(event *Event, subscribers SubscriberList) {
	for _, subscriber := range subscribers {
		if event.CreationDate > subscriber.creationDate {

			count := 0
			if value, isOk := subscriber.Counts.Load(event.Name); isOk {
				count = value.(int)
			}

			if subscriber.Repeat == 0 || subscriber.Repeat > count {
				subscriber.handler(event)

				count += 1

				subscriber.Counts.Store(event.Name, count)

				if event.isCompleted {
					break
				}
			}
		}
	}
}
