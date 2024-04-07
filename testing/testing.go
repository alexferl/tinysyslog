package testing

import (
	"tinysyslog/config"
)

func init() {
	c := config.New()
	c.BindFlags()
}
