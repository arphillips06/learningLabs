package model

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RoutingOptions struct {
	RouterID string `yaml:"router-id"`
	ASN      int    `yaml:"asn"`
}

type Inventory struct {
	Defaults InventoryDefaults    `yaml:"defaults"`
	Devices  map[string]InvDevice `yaml:"devices"`
}

type InventoryDefaults struct {
	NetconfPort int    `yaml:"netconf_port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	LicenseFile string `yaml:"license_file"`
}

type InvDevice struct {
	Address string `yaml:"address"`
	Vars    string `yaml:"vars"`
}

func LoadInventory(path string) (*Inventory, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inv Inventory
	return &inv, yaml.Unmarshal(b, &inv)
}
