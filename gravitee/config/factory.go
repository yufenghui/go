package config

import "github.com/jinzhu/configor"

var DefaultConfig = &Config{
	Database: DatabaseConfig{
		Type:         "mysql",
		Host:         "localhost",
		Port:         3306,
		User:         "root",
		Password:     "root",
		DatabaseName: "gravitee",
		MaxIdleConns: 5,
		MaxOpenConns: 5,
	},
	Oauth: OauthConfig{
		AccessTokenLifetime:  3600,    // 1 hour
		RefreshTokenLifetime: 1209600, // 14 days
		AuthCodeLifetime:     3600,    // 1 hour
	},
	Session: SessionConfig{
		Secret:   "test_secret",
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
	},
	IsDevelop: true,
}

func NewDefaultConfig() *Config {
	return DefaultConfig
}

func NewConfig(configFile string) *Config {
	if configFile != "" {
		config := &Config{}
		configor.Load(config, configFile)
		return config
	}

	return DefaultConfig
}
