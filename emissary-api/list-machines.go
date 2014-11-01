package emissaryapi

import ()

type Machine struct {
	Name    string
	Address string
}

func (c *ApiClient) ListMachines() (machines []Machine, err error) {
	nodes, _, err := c.consul.Catalog().Nodes(&c.q)
	if err != nil {
		return
	}
	machines = make([]Machine, len(nodes))
	for i, v := range nodes {
		machines[i] = Machine{v.Node, v.Address}
	}
	return
}
func (c *ApiClient) ListMachinesPattern(patterns ...string) (machines []Machine, err error) {
	nodes, _, err := c.consul.Catalog().Nodes(&c.q)
	if err != nil {
		return
	}
	machines = make([]Machine, 0, len(nodes))
	for _, v := range nodes {
		if matchAny(patterns, v.Node) {
			machines = append(machines, Machine{v.Node, v.Address})
		}
	}
	return
}
