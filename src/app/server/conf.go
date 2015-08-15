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

production:
  debug: false
  port: 5000
  title: lmbd
  api:
    prefix: /api
`)
	Must(err)
	conf = c
}
