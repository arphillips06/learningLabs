package model

type DeviceConfig struct {
	System         System          `yaml:"system"`
	Loopback       Loopback        `yaml:"loopback"`
	CoreInterfaces []CoreInterface `yaml:"core_interfaces"`
	Protocols      Protocols       `yaml:"protocols"`
	RoutingOptions *RoutingOptions `yaml:"routing-options,omitempty"`
}

type System struct {
	Hostname     string `yaml:"hostname"`
	RootPassword string `yaml:"root_password"`
}

type Loopback struct {
	Lo0 []LoUnit `yaml:"lo0"`
}

type LoUnit struct {
	Unit int    `yaml:"unit"`
	IPv4 string `yaml:"ipv4"`
	IPv6 string `yaml:"ipv6"`
	ISO  string `yaml:"iso"`
}

type CoreInterface struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	// Unit optional: nil means "no unit stanza"
	Unit *int   `yaml:"unit"`
	IPv4 string `yaml:"ipv4"`

	// Optional LAG membership for ge-* physicals
	AeBundle string `yaml:"ae_bundle"` // e.g. "ae0"
}
