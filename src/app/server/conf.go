package server

import (
	. "app/server/utils"

	"github.com/olebedev/config"
)

var conf *config.Config

func init() {
	c, err := config.ParseYaml(`

local:
  debug: true
  port: 5000
  title: lmbd
  db: ./db.sqlite
  api:
    prefix: /api
  duktape:
    pool:
      use: false

production:
  debug: false
  port: 5000
  title: lmbd
  api:
    prefix: /api
  duktape:
    pool:
      use: true
      size: 1
`)
	Must(err)
	conf = c
}
