package emissaryapi

import (
	"time"
)

type Event struct {
	Name string
	Data []byte
}

func (c *ApiClient) EventListener(name string, interval time.Duration) chan Event {
	ch := make(chan Event)
	go func() {
		prevEvents := make(map[string]bool, 256)
		for {
			events, _, err := c.consul.Event().List(name, &c.q)
			if err == nil {
				for _, v := range events {
					if !prevEvents[v.ID] {
						ch <- Event{v.Name, v.Payload}
					}
					prevEvents[v.ID] = true
				}
			}
			prevEvents = make(map[string]bool, 256)
			time.Sleep(interval)
		}
	}()
	return ch
}
