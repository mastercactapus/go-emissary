package emissaryapi

import (
	"github.com/armon/consul-api"
)

type ApiClient struct {
	consul *consulapi.Client
	dc     string
	sess   string
	kv     consulapi.KV
	q      consulapi.QueryOptions
	w      consulapi.WriteOptions
}
type kvLock struct {
	session string
	c       *ApiClient
}

func NewClient(c *consulapi.Client, datacenter string) *ApiClient {
	return &ApiClient{
		consul: c,
		dc:     datacenter,
		q:      consulapi.QueryOptions{Datacenter: datacenter},
		w:      consulapi.WriteOptions{Datacenter: datacenter},
		kv:     *c.KV(),
	}
}

func (c *ApiClient) kvGet(key string) ([]byte, error) {
	pair, _, err := c.kv.Get(key, &c.q)
	if err != nil {
		return nil, err
	}
	return pair.Value, nil
}
func (c *ApiClient) kvSet(key string, value []byte) error {
	_, err := c.kv.Put(&consulapi.KVPair{Key: key, Value: value}, &c.w)
	return err
}
func (c *ApiClient) kvSetSession(key string, value []byte) error {
	if c.sess != "" {
		return ErrNoSession
	}
	_, err := c.kv.Put(&consulapi.KVPair{Key: key, Value: value, Session: c.sess}, &c.w)
	return err
}
