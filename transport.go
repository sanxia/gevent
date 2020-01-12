package gevent

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ITransport interface {
		Load() *Event
		Store(*Event)
	}

	defaultTransport struct {
		eventChan chan *Event //event buffer channel
		backCount int         //maximum backup buffer
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize the default transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewDefaultTransport(backCount int) ITransport {
	return &defaultTransport{
		eventChan: make(chan *Event, backCount),
		backCount: backCount,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get event data from transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultTransport) Load() *Event {
	var event *Event

	select {
	case eventSource := <-s.eventChan:
		event = eventSource
	}

	return event
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set event data to transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultTransport) Store(event *Event) {
	s.eventChan <- event
}
