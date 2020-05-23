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
	jsonBody, _ := glib.ToJson(event)
	return jsonBody
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * deserialize json
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *defaultSerialization) Deserialize(jsonBody string) *Event {
	var event *Event

	if len(jsonBody) > 0 {
		glib.FromJson(jsonBody, &event)
	}

	return event
}
