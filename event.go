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
		Name         string      //event name
		Data         interface{} //event data
		IsBroadcast  bool        //it is broadcast
		isCompleted  bool        //it is complete
		CreationDate int64       //event creation time
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * mapping event data
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Event) Map(mapFunc func(data interface{}) interface{}) {
	s.Data = mapFunc(s.Data)
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
