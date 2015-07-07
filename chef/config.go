package chef

import (
	chefGo "github.com/go-chef/chef"
	"io/ioutil"
)

type Config struct {
	Name    string
	Key     string
	BaseURL string
}

func (c *Config) Client() (*chefGo.Client, error) {
	key, err := ioutil.ReadFile(c.Key)
	if err != nil {
		return nil, err
	}
	config := chefGo.Config{
		Name:    c.Name,
		Key:     string(key),
		BaseURL: c.BaseURL,
		SkipSSL: false,
	}

	return chefGo.NewClient(&config)
}
