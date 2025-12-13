package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/Juniper/go-netconf/netconf"
	"gopkg.in/yaml.v3"

	"learninglabs/sp-6node/internal/junos"
	"learninglabs/sp-6node/internal/model"
)

func main() {
	// ------------------------------------------------------------------
	// Load inventory
	// ------------------------------------------------------------------
	inv, err := model.LoadInventory("inventory/inventory.yml")
	if err != nil {
		log.Fatalf("failed to load inventory: %v", err)
	}

	// ------------------------------------------------------------------
	// Parse templates once
	// ------------------------------------------------------------------
	tmpl, err := template.ParseGlob("templates/*.xml.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.ParseGlob("templates/protocols/*.xml.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------------
	// Iterate devices with per-device timeout
	// ------------------------------------------------------------------
	for name, dev := range inv.Devices {
		log.Printf("===== %s =====", name)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := deployWithTimeout(ctx, name, dev, inv, tmpl)
		cancel()

		if err != nil {
			log.Printf("[%s] ERROR: %v", name, err)
			continue
		}
	}

	log.Println("All devices processed")
}

func deployWithTimeout(
	ctx context.Context,
	name string,
	dev model.InvDevice,
	inv *model.Inventory,
	tmpl *template.Template,
) error {

	done := make(chan error, 1)

	go func() {
		done <- deployDevice(name, dev, inv, tmpl)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout after 10s")
	case err := <-done:
		return err
	}
}

func deployDevice(
	name string,
	dev model.InvDevice,
	inv *model.Inventory,
	tmpl *template.Template,
) error {

	// --------------------------------------------------------------
	// Load device vars
	// --------------------------------------------------------------
	data, err := os.ReadFile(dev.Vars)
	if err != nil {
		return err
	}

	var cfg model.DeviceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	role := model.DetectRole(&cfg)
	log.Printf("[%s] detected role: %s", name, role)

	// --------------------------------------------------------------
	// Render configuration
	// --------------------------------------------------------------
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "configuration", cfg); err != nil {
		return err
	}

	// --------------------------------------------------------------
	// NETCONF connect
	// --------------------------------------------------------------
	addr := fmt.Sprintf(
		"%s:%d",
		dev.Address,
		inv.Defaults.NetconfPort,
	)

	session, err := netconf.DialSSH(
		addr,
		netconf.SSHConfigPassword(
			inv.Defaults.Username,
			inv.Defaults.Password,
		),
	)
	if err != nil {
		return err
	}
	defer session.Close()

	// Upload + install license on all devices
	if inv.Defaults.LicenseFile != "" {
		log.Printf("[%s] installing license", name)

		remote := "/var/tmp/junos.lic"

		if err := junos.SCPUpload(
			inv.Defaults.LicenseFile,
			inv.Defaults.Username,
			dev.Address,
			remote,
		); err != nil {
			return err
		}

		if err := junos.InstallLicense(session, remote); err != nil {
			return err
		}
	}

	// --------------------------------------------------------------
	// Lock candidate
	// --------------------------------------------------------------
	if _, err := session.Exec(netconf.RawMethod(`
<lock>
  <target><candidate/></target>
</lock>`)); err != nil {
		return err
	}
	defer session.Exec(netconf.RawMethod(`
<unlock>
  <target><candidate/></target>
</unlock>`))

	// --------------------------------------------------------------
	// Load configuration (merge)
	// --------------------------------------------------------------
	rpc := fmt.Sprintf(`
<load-configuration action="merge" format="xml">
%s
</load-configuration>`, buf.String())

	if _, err := session.Exec(netconf.RawMethod(rpc)); err != nil {
		return err
	}

	// --------------------------------------------------------------
	// Commit check
	// --------------------------------------------------------------
	if _, err := session.Exec(netconf.RawMethod(`
<commit-configuration check="true"/>
`)); err != nil {
		return err
	}

	// --------------------------------------------------------------
	// Commit
	// --------------------------------------------------------------
	if _, err := session.Exec(netconf.RawMethod(`
<commit-configuration/>
`)); err != nil {
		return err
	}

	log.Printf("[%s] commit successful", name)
	return nil
}
