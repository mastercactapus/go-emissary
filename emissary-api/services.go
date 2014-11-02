package emissaryapi

import (
	"github.com/armon/consul-api"
	"time"
)

func (c *ApiClient) ServicePass(name, note string) error {
	return c.consul.Agent().PassTTL("service:emissary:"+name, note)
}
func (c *ApiClient) ServiceFail(name, note string) error {
	return c.consul.Agent().FailTTL("service:emissary:"+name, note)
}

func (c *ApiClient) RegisterService(unitName string, monitorTtl time.Duration) error {
	service := &consulapi.AgentServiceRegistration{Name: unitName, ID: "emissary:" + unitName}
	if monitorTtl > 0 {
		service.Check = &consulapi.AgentServiceCheck{
			TTL: monitorTtl.String(),
		}
	}
	return c.consul.Agent().ServiceRegister(service)
}

func (c *ApiClient) DeregisterService(unitName string) error {
	return c.consul.Agent().ServiceDeregister("emissary:" + unitName)
}
