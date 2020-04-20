package gevent

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ISerialization interface {
		Serialize(event *Event) string
		Deserialize(body string) *Event
	}

	defaultSerialization struct {
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * initialize the default serialization
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewDefaultSerialization() ISerialization {
	return &defaultSerialization{}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * serialize json
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultSerialization) Serialize(event *Event) string {
	eventJson, _ := glib.ToJson(event)
	return eventJson
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * deserialize json
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultSerialization) Deserialize(body string) *Event {
	var event *Event

	if len(body) > 0 {
		glib.FromJson(body, &event)
	}

	return event
}
