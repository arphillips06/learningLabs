# Junos Automation Tool

Minimal Go-based tool to render and deploy Junos configuration using **NETCONF** and **SCP**.

- Renders Junos XML from Go templates + YAML intent
- Deploys config via NETCONF (candidate → commit)
- Uploads license files using system `scp`
- Installs licenses with `request system license add`
- Detects P vs PE automatically
- Skips unreachable devices after a 10s timeout
