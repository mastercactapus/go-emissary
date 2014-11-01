package emissaryapi

import (
	"github.com/armon/consul-api"
)

type ApiClient struct {
	Store  UnitStore
	consul *consulapi.Client
	dc     string
	q      consulapi.QueryOptions
}

func NewClient(c *consulapi.Client, datacenter string) *ApiClient {
	return &ApiClient{
		consul: c,
		dc:     datacenter,
		Store:  *NewUnitStore(c, datacenter),
		q:      consulapi.QueryOptions{Datacenter: datacenter},
	}
}
