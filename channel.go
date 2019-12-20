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
	Channel struct {
		Name        string                    //Channel name
		eventChan   chan *Event               //event buffer channel
		subscribers map[string]SubscriberList //subscriber Collection
		mu          sync.Mutex                //synchronous
		isRunning   bool                      //is run
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * publishing events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Publish(eventName string, data interface{}) *Channel {
	event := &Event{
		Name:         eventName,
		Data:         data,
		CreationDate: time.Now().UnixNano(),
	}

	go func() {
		s.dispatchEvent(event)
	}()

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * subscribe to events
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) Subscribe(eventName string, eventHandler EventHandler, args ...int) *Channel {
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
func (s *Channel) Unsubscribe(eventNames ...string) *Channel {
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
func (s *Channel) dispatchEvent(event *Event) *Channel {
	if _, isOk := s.subscribers[event.Name]; isOk {
		s.eventChan <- event
	}

	return s
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event loop
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Channel) eventLoop() {
	for {
		select {
		case event := <-s.eventChan:
			for eventName, subscribers := range s.subscribers {
				if eventName == event.Name {
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
	for _, subscriber := range subscribers {
		if event.CreationDate > subscriber.creationDate {

			s.mu.Lock()

			subscriber.Counts[event.Name] += 1

			s.mu.Unlock()

			if subscriber.Repeat == 0 || subscriber.Counts[event.Name] <= subscriber.Repeat {
				subscriber.handler(event)

				if event.isCompleted {
					break
				}
			}

		}
	}
}
