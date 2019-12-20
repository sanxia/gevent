package gevent

/* ================================================================================
 * gevent
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	SubscriberList []*Subscriber
	Subscriber     struct {
		Priority     int            //subscription priority
		Repeat       int            //Number of triggers（0:Infinite | 1:once | ...）
		Counts       map[string]int //event counter
		handler      EventHandler   //event processor
		creationDate int64          //subscription time (nanoseconds)
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * sort Interface - collection length
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (list SubscriberList) Len() int {
	return len(list)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * sort interface - collection item comparison
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (list SubscriberList) Less(i, j int) bool {
	if list[i].Priority < list[j].Priority {
		return true
	}

	return false
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * sort interface - swap collection items
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (list SubscriberList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
