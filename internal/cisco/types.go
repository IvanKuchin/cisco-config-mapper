package cisco

type Iface struct {
	Name    string
	Content []string
}

type Cisco struct {
	Before_ifaces []string
	Ifaces        []Iface
	After_ifaces  []string
}
