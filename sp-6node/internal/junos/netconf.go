package junos

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

type Client struct {
	s *netconf.Session
}

func Dial(addr, user, pass string) (*Client, error) {
	s, err := netconf.DialSSH(addr, netconf.SSHConfigPassword(user, pass))
	if err != nil {
		return nil, err
	}
	return &Client{s: s}, nil
}

func (c *Client) Close() error { return c.s.Close() }

func (c *Client) exec(xml string) (string, error) {
	reply, err := c.s.Exec(netconf.RawMethod(xml))
	if err != nil {
		return "", err
	}
	return reply.Data, nil
}

func (c *Client) LockCandidate() error {
	_, err := c.exec(`<lock><target><candidate/></target></lock>`)
	return err
}

func (c *Client) UnlockCandidate() error {
	_, err := c.exec(`<unlock><target><candidate/></target></unlock>`)
	return err
}

func (c *Client) LoadConfigXML(merge bool, xmlConfig string) error {
	action := "merge"
	if !merge {
		action = "replace"
	}
	rpc := fmt.Sprintf(`
<load-configuration action="%s" format="xml">
  %s
</load-configuration>`, action, xmlConfig)
	_, err := c.exec(rpc)
	return err
}

func (c *Client) CommitCheck() error {
	_, err := c.exec(`<commit-configuration check="true"/>`)
	return err
}

func (c *Client) Commit() error {
	_, err := c.exec(`<commit-configuration/>`)
	return err
}
