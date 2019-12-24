package gevent

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type EventHandler func(*Event)
type (
	EventList []*Event
	Event     struct {
		ChannelName  string      `json:"channel_name"`  //channel name
		Name         string      `json:"name"`          //event name
		Data         interface{} `json:"data"`          //event data
		IsBroadcast  bool        `json:"is_broadcast"`  //it is broadcast
		isCompleted  bool        `json:"_"`             //it is complete
		CreationDate int64       `json:"creation_date"` //event creation time
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * mapping event data
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Event) Map(mapFunc func(data interface{}) interface{}) {
	if !s.IsBroadcast {
		s.Data = mapFunc(s.Data)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * the next subscriber of the event continues to process
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Event) Next() {
	if s.isCompleted {
		return
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * event Completed
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Event) Completed() {
	s.isCompleted = true
}
