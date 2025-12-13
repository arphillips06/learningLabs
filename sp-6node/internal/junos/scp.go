package junos

import (
	"fmt"
	"os/exec"
)

func SCPUpload(
	localPath string,
	username string,
	host string,
	remotePath string,
) error {

	target := fmt.Sprintf("%s@%s:%s", username, host, remotePath)

	cmd := exec.Command(
		"scp",
		"-q",
		localPath,
		target,
	)

	// Inherit environment (SSH_AUTH_SOCK, PATH, etc.)
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}
