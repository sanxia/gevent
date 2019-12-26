package gevent

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IEventSource interface {
		GetChannelName() string
		GetEventName() string
		GetData() interface{}
	}
)
