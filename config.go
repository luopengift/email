package email

import "github.com/luopengift/types"

// Config config
type Config struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}

// Parse smtp from v
func (c *Config) Parse(v interface{}) error {
	return types.Format(v, c)
}
