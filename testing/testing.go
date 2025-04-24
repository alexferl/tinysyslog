package testing

import (
	"github.com/alexferl/tinysyslog/config"
)

func init() {
	c := config.New()
	c.BindFlags()
}
