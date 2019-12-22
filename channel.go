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
		Name        string                    //Channel name
		store       IStore                    //event store
		subscribers map[string]SubscriberList //subscriber Collection
		mu          sync.Mutex                //synchronous
	}
)

func NewChannel(channelName string, store IStore) IChannel {
	channel := &Channel{
		Name:        channelName,
		store:       store,
		subscribers: make(map[string]SubscriberList, 0),
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
	s.mu.Lock()
	defer s.mu.Unlock()

	subscribers, isOk := s.subscribers[eventName]
	if !isOk {
		subscribers = make(SubscriberList, 0)
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
		Counts:       make(map[string]int, 0),
		handler:      eventHandler,
		creationDate: time.Now().UnixNano(),
	})

	//优先级从高到低降序
	sort.Sort(sort.Reverse(subscribers))

	s.subscribers[eventName] = subscribers

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * unsubscribe events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Unsubscribe(eventNames ...string) IChannel {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(eventNames) == 0 {
		s.subscribers = make(map[string]SubscriberList)
	} else {
		for _, eventName := range eventNames {
			if _, isOk := s.subscribers[eventName]; isOk {
				s.subscribers[eventName] = make(SubscriberList, 0)
			}
		}
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * dispatch event
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) dispatchEvent(event *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.IsBroadcast {
		s.store.Put(event)
	} else {
		if _, isOk := s.subscribers[event.Name]; isOk {
			s.store.Put(event)
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event loop
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) eventLoop() {
	for {
		if event := s.store.Get(); event != nil {
			for eventName, subscribers := range s.subscribers {
				if eventName == event.Name {
					go s.eventHandler(event, subscribers)
				} else if event.IsBroadcast {
					go s.eventHandler(event, subscribers)
				}
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event handling
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) eventHandler(event *Event, subscribers SubscriberList) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, subscriber := range subscribers {
		if event.CreationDate > subscriber.creationDate {

			subscriber.Counts[event.Name] += 1

			if subscriber.Repeat == 0 || subscriber.Counts[event.Name] <= subscriber.Repeat {
				subscriber.handler(event)

				if event.isCompleted {
					break
				}
			}
		}
	}
}
