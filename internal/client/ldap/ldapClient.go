package ldap

import (
	ld "github.com/go-ldap/ldap"
	"github.com/phluxx/gogert/internal/service/config"
)

type Client struct {
	ldap *ld.Conn
}

func New(cfg *config.LdapConfig) (*Client, error) {
	return &Client{}, nil
	l, err := ld.Dial("tcp", cfg.LdapHost+":"+cfg.LdapPort)
	if err != nil {
		return nil, err
	}

	c := &Client{
		ldap: l,
	}
	return c, nil
}

func (c *Client) Connect() error {
	return nil
}
func (c *Client) Authenticate(username, password string) error {
	return nil
}

func (c *Client) ChangePassword(username, newPassword string) error {
	return nil
}
