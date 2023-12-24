package config

type Iface_mapping struct {
	From    string `yaml:"from"`
	To      string `yaml:"to"`
	Visited bool
}

type Config struct {
	Iface_mappings  []Iface_mapping `yaml:"interface-mappings"`
	Remove_prefixes []string        `yaml:"remove-prefixes"`
	Remove_lines    []string        `yaml:"remove-lines"`
	Append          string          `yaml:"append"`
}
