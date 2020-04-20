package gevent

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ITransport interface {
		Load(callback func(*Event) error) error
		Store(*Event) error
	}

	defaultTransport struct {
		eventChan      chan *Event //event buffer channel
		maxBufferCount int         //maximum backup buffer
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize the default transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewDefaultTransport(maxBufferCount int) ITransport {
	if maxBufferCount <= 0 {
		maxBufferCount = 64
	}

	return &defaultTransport{
		eventChan:      make(chan *Event, maxBufferCount),
		maxBufferCount: maxBufferCount,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get event data from transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultTransport) Load(callback func(*Event) error) error {
	for eventSource := range s.eventChan {
		if callback != nil {
			callback(eventSource)
		}

		return nil
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set event data to transport
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultTransport) Store(event *Event) error {
	s.eventChan <- event

	return nil
}
