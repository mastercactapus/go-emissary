package emissaryapi

import (
	"time"
)

type Event struct {
	Name    string
	Payload []byte
}

func (c *ApiClient) EventListener(eventName string, interval time.Duration) (events chan *Event, stop chan int) {
	events = make(chan *Event, 256)
	stop = make(chan int)
	ticker := time.NewTicker(interval)

	go func() {
		firedEvents := make(map[string]bool, 1024)
		numEvents := 0
		for {
			select {
			case <-ticker.C:
				e, _, err := c.consul.Event().List(eventName, &c.q)
				if err != nil {
					continue
				}
				var update map[string]bool
				if numEvents < 1024 {
					update = firedEvents
				} else {
					update = make(map[string]bool, 1024)
				}
				for _, v := range e {
					if !firedEvents[v.Name] {
						events <- &Event{v.Name, v.Payload}
					}
					update[v.Name] = true
				}
				firedEvents = update
			case <-stop:
				ticker.Stop()
				close(events)
				return
			}
		}
	}()

	return
}
