# gevent
golang event hubs publishing and subscriptions

=====
Example
=====

    import "time"

    import "fmt"

    import "log"

    import "git.hub.com/sanxia/gevent"


    eventHub := gevent.NewEventHub(10)

    channel := eventHub.GetChannel("test")

    channel.Publish("article_deleted", fmt.Sprintf("%s", "yichang"))

    channel.Subscribe("user_created", func(e *gevent.Event) {

        log.Printf("1 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Map(func(data interface{}) interface{} {

            return "world"

        })

        e.Next()

    }, 9, 1).Subscribe("user_created", func(e *gevent.Event) {

        log.Printf("2 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Next()

    }, 1, 1).Subscribe("user_created", func(e *gevent.Event) {

        log.Printf("3 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Next()

    }, 3)

    channel.Subscribe("article_deleted", func(e *gevent.Event) {

        log.Printf("4 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Next()

    }, 3)

    channel.Subscribe("city_changed", func(e *gevent.Event) {

        log.Printf("5 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Next()

    }).Unsubscribe("city_changed")

    channel.Subscribe("system_notify", func(e *gevent.Event) {

        log.Printf("6 event: %s, data: %s, creation: %d",
        e.Name, e.Data, e.CreationDate)

        e.Next()

    })

    eventHub.Broadcast("system_notify", fmt.Sprintf("%s", "broadcast"))

    for i := 0; i < 2; i++ {

        channel.Publish("user_created", fmt.Sprintf("%s", "hello")).Publish("city_changed", fmt.Sprintf("%s", "beijing"))

        time.Sleep(600 * time.Millisecond)

    }

-----
Run Result Response
-----

    6 event: system_notify, data: broadcast, creation: 1576857124739584000

    1 event: user_created, data: hello, creation: 1576857124739562000

    3 event: user_created, data: world, creation: 1576857124739562000

    2 event: user_created, data: world, creation: 1576857124739562000

    3 event: user_created, data: hello, creation: 1576857125343876000

