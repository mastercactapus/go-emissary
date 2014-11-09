package emissaryapi

import (
	"errors"
	"github.com/armon/consul-api"
	"time"
)

type Machine struct {
	Name       string
	Datacenter string
	Address    string
	Domain     string
	Metadata   []string
}

var ErrNoSession = errors.New("Cannot perform action without a session. Call RegisterSelf first.")
var ErrAlreadyRegistered = errors.New("Machine is already registered")
var stopHeartbeat = make(chan bool)

func (c *ApiClient) RegisterSelf(meta []string, ttl time.Duration) error {
	if c.sess != "" {
		return ErrAlreadyRegistered
	}
	var se consulapi.SessionEntry
	se.Checks = []string{"serfHealth", "service:emissary"}
	var service = &consulapi.AgentServiceRegistration{
		ID:    "emissary",
		Name:  "emissary",
		Tags:  meta,
		Check: &consulapi.AgentServiceCheck{TTL: ttl.String()},
	}
	err := c.consul.Agent().ServiceRegister(service)
	if err != nil {
		return err
	}
	err = c.consul.Agent().PassTTL("service:emissary", "")
	if err != nil {
		return err
	}
	se.Name = "emissary:emissaryd"
	id, _, err := c.consul.Session().Create(&se, &c.w)
	if err != nil {
		return err
	}
	c.sess = id

	go func() {
		ticker := time.NewTicker(ttl / 2)
		for {
			select {
			case <-ticker.C:
				c.consul.Agent().PassTTL("service:emissary", "")
			case <-stopHeartbeat:
				ticker.Stop()
				break
			}
		}
	}()

	return nil
}
func (c *ApiClient) UnregisterSelf() error {
	if c.sess == "" {
		return ErrNoSession
	}
	_, err := c.consul.Session().Destroy(c.sess, &c.w)
	if err != nil {
		return err
	}
	c.sess = ""
	stopHeartbeat <- true
	return c.consul.Agent().ServiceDeregister("emissary")
}

func (c *ApiClient) Self() (*Machine, error) {
	name, err := c.consul.Agent().NodeName()
	if err != nil {
		return nil, err
	}
	self, err := c.consul.Agent().Self()
	if err != nil {
		return nil, err
	}
	svcs, err := c.consul.Agent().Services()
	if err != nil {
		return nil, err
	}
	emSvc := svcs["emissary"]

	m := new(Machine)
	m.Name = name
	m.Address = self["Config"]["AdvertiseAddr"].(string)
	m.Datacenter = self["Config"]["Datacenter"].(string)
	m.Domain = self["Config"]["Domain"].(string)
	if emSvc != nil {
		m.Metadata = emSvc.Tags
	}
	return m, nil
}

func (m *Machine) FQDN() string {
	return m.Name + ".node." + m.Datacenter + "." + m.Domain
}

func (c *ApiClient) Machines(datacenter, meta string) ([]Machine, error) {
	var opts consulapi.QueryOptions
	opts.Datacenter = datacenter
	self, err := c.Self()
	if err != nil {
		return nil, err
	}
	s, _, err := c.consul.Catalog().Service("emissary", meta, &opts)
	if err != nil {
		return nil, err
	}
	machines := make([]Machine, 0, len(s))
	for _, v := range s {
		var m Machine
		m.Address = v.Address
		m.Name = v.Node
		if datacenter == "" {
			m.Datacenter = self.Datacenter
		} else {
			m.Datacenter = datacenter
		}
		m.Metadata = v.ServiceTags
		m.Domain = self.Domain
		machines = append(machines, m)
	}
	return machines, nil
}
