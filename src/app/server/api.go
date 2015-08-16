package server

import (
	"github.com/gin-gonic/gin"
)

// Api is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type Api struct{}

// Bind attaches api routes
func (api *Api) Bind(group *gin.RouterGroup) {
	group.GET("/v1/conf", api.ConfHandler)
}

// Serve the app config, for example
func (_ *Api) ConfHandler(c *gin.Context) {
	app := c.MustGet("app").(*App)
	c.JSON(200, app.Conf.Root)
}
