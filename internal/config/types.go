package config

type Iface_mapping struct {
	From            string   `yaml:"from"`
	To              string   `yaml:"to"`
	Append          string   `yaml:"append"`
	Prepend         string   `yaml:"prepend"`
	Remove_prefixes []string `yaml:"remove-prefixes"`
	Remove_lines    []string `yaml:"remove-lines"`
	Visited         bool
}

type Config struct {
	Prepend         string          `yaml:"prepend"`
	Iface_mappings  []Iface_mapping `yaml:"interface-mappings"`
	Remove_prefixes []string        `yaml:"remove-prefixes"`
	Remove_lines    []string        `yaml:"remove-lines"`
	Append          string          `yaml:"append"`
}

func (iface_mapping *Iface_mapping) MarkVisited() {
	iface_mapping.Visited = true
}
