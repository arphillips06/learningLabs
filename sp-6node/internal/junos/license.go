package junos

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

func InstallLicense(session *netconf.Session, filename string) error {
	rpc := fmt.Sprintf(
		`<request-license-add xmlns="http://xml.juniper.net/junos/*/junos">`+
			`<filename>%s</filename>`+
			`</request-license-add>`,
		filename,
	)

	_, err := session.Exec(netconf.RawMethod(rpc))
	return err
}
