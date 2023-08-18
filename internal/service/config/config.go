package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HttpConfig
	LdapConfig
	JwtConfig
}

type HttpConfig struct {
	// HTTP
	HttpPort string `envconfig:"HTTP_PORT" default:"8080"`
}

type LdapConfig struct {
	// LDAP
	LdapHost string `envconfig:"LDAP_HOST" default:"localhost"`
	LdapPort string `envconfig:"LDAP_PORT" default:"389"`
	LdapBase string `envconfig:"LDAP_BASE" default:"dc=ewnix,dc=net"`
	LdapBind string `envconfig:"LDAP_BIND" default:"cn=admin,dc=ewnix,dc=net"`
	LdapPass string `envconfig:"LDAP_PASS" default:"password"`
}

type JwtConfig struct {
	// JWT
	Secret string `envconfig:"JWT_SECRET" default:"fellysmart"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
