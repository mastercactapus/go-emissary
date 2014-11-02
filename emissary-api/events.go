package emissaryapi

import (
	"time"
)

func (c *ApiClient) EventListener(name string, interval time.Duration) chan []byte {
	ch := make(chan []byte)
	go func() {
		prevEvents := make(map[string]bool, 256)
		for {
			events, _, err := c.consul.Event().List(name, &c.q)
			if err == nil {
				for _, v := range events {
					if !prevEvents[v.ID] {
						ch <- v.Payload
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
