package model

type Role string

const (
	RoleP  Role = "P"
	RolePE Role = "PE"
)

// DetectRole determines router role from config intent
// Rule: PE routers have BGP, P routers do not
func DetectRole(cfg *DeviceConfig) Role {
	if cfg.Protocols.BGP != nil {
		return RolePE
	}
	return RoleP
}
