package emissaryapi

import (
	"github.com/armon/consul-api"
)

type ApiClient struct {
	Store  UnitStore
	consul *consulapi.Client
	dc     string
	kv     consulapi.KV
	q      consulapi.QueryOptions
	w      consulapi.WriteOptions
}

func NewClient(c *consulapi.Client, datacenter string) *ApiClient {
	return &ApiClient{
		consul: c,
		dc:     datacenter,
		Store:  *NewUnitStore(c, datacenter),
		q:      consulapi.QueryOptions{Datacenter: datacenter},
		w:      consulapi.WriteOptions{Datacenter: datacenter},
		kv:     *c.KV(),
	}
}
