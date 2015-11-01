package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// API is a defined as struct bundle
// for api. Feel free to organize
// your app as you wish.
type API struct{}

// Bind attaches api routes
func (api *API) Bind(group *gin.RouterGroup) {
	group.GET("/v1/conf", api.ConfHandler)
}

// ConfHandler handle the app config, for example
func (api *API) ConfHandler(c *gin.Context) {
	app := c.MustGet("app").(*App)
	<-time.After(time.Millisecond * 500)
	c.JSON(200, app.Conf.Root)
}
