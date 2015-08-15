package utils

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/olebedev/config"
)

// Must raises an error if it not nil
func Must(e error) {
	if e != nil {
		panic(e)
	}
}

type Kit struct {
	Conf   *config.Config
	Engine *gin.Engine
}

func NewKit(c *cli.Context, conf *config.Config) *Kit {
	co, err := conf.Get(c.String("config"))
	// parse environ variables
	co.Env()
	Must(err)

	// set up gin
	if !co.UBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Kit{
		Conf: co,
	}
}
