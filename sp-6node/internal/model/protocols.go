package model

type Protocols struct {
	BGP  *BGP `yaml:"bgp,omitempty"`
	ISIS ISIS `yaml:"isis"`
	MPLS MPLS `yaml:"mpls"`
	LDP  LDP  `yaml:"ldp"`
}

type ISIS struct {
	Metric          int      `yaml:"metric"`
	NodeSIDIPv4     int      `yaml:"node_sid_ipv4"`
	NodeSIDIPv6     int      `yaml:"node_sid_ipv6"`
	P2PInterfaces   []string `yaml:"p2p_interfaces"`
	PassiveLoopback bool     `yaml:"passive_loopback"`
}

type MPLS struct {
	Interfaces []string `yaml:"interfaces"`
}

type LDP struct {
	Interfaces []string `yaml:"interfaces"`
	P2MP       bool     `yaml:"p2mp"`
}

type BGP struct {
	Group        string        `yaml:"group"`
	Type         string        `yaml:"type"`
	LocalAddress string        `yaml:"local-address"`
	Families     []BGPFamily   `yaml:"families"`
	Neighbors    []BGPNeighbor `yaml:"neighbors"`
}

type BGPFamily struct {
	Name string `yaml:"name"`
	NLRI string `yaml:"NLRI"`
}

type BGPNeighbor struct {
	Address     string `yaml:"address"`
	Description string `yaml:"description"`
}
