package server

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/config"
)

// App struct.
// There is no singleton anti-pattern,
// all variables defined locally inside
// this struct.
type App struct {
	Engine *gin.Engine
	Conf   *config.Config
	React  *React
	API    *API
}

// NewApp returns initialized struct
// of main server application.
func NewApp(opts ...AppOptions) *App {
	options := AppOptions{}
	for _, i := range opts {
		options = i
		break
	}
	options.init()

	// Parse config yaml string from ./conf.go
	conf, err := config.ParseYaml(confString)
	Must(err)
	// Choise a config section by given string
	conf, err = conf.Get(options.Config)
	Must(err)

	// Parse environ variables for defined
	// in config constants
	conf.Env()

	// Set up gin
	if !conf.UBool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Make an engine
	engine := gin.Default()

	// Initialize the application
	app := &App{
		Conf:   conf,
		Engine: engine,
		API:    &API{},
		React: NewReact(
			conf.UString("duktape.path"),
			conf.UBool("debug"),
			engine,
		),
	}

	// Define routes and middlewares
	app.Engine.StaticFS("/static", &assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "static",
	})

	// Map app struct to access from request handlers
	// and middlewares
	app.Engine.Use(func(c *gin.Context) {
		c.Set("app", app)
	})

	// Avoid favicon react handling
	app.Engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(301, "/static/images/favicon.ico")
	})

	// Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(
			app.Conf.UString("api.prefix"),
		),
	)

	// Map uuid for every requests
	app.Engine.Use(func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
	})

	// Handle all not found routes via react app
	app.Engine.NoRoute(app.React.Handle)

	return app
}

// Run runs the app
func (app *App) Run() {
	Must(app.Engine.Run(":" + app.Conf.UString("port")))
}

// AppOptions is options struct
type AppOptions struct {
	Config string
}

func (ao *AppOptions) init() {
	if ao.Config == "" {
		ao.Config = "local"
	}
}
