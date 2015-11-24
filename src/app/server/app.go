package server

import (
	"io"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/itsjamie/go-bindata-templates"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/config"
)

// App struct.
// There is no singleton anti-pattern,
// all variables defined locally inside
// this struct.
type App struct {
	Engine *echo.Echo
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

	// Make an engine
	engine := echo.New()

	// Set up echo
	engine.SetDebug(conf.UBool("debug"))

	// Regular middlewares
	engine.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

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

	// Load embedded templates
	app.Engine.SetRenderer(echoRenderer{})

	// Map app struct to access from request handlers
	// and middlewares
	app.Engine.Use(func(c *echo.Context) error {
		c.Set("app", app)
		return nil
	})

	// Map uuid for every requests
	app.Engine.Use(func(c *echo.Context) error {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
		return nil
	})

	// Create file http server from bindata
	fileServerHandler := http.FileServer(&assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	})

	app.Engine.Get("/static/*", func(c *echo.Context) error {
		if _, err := Asset(c.Request().URL.Path[1:]); err == nil {
			fileServerHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
		return echo.NewHTTPError(http.StatusNotFound)
	})

	// Avoid favicon react handling
	app.Engine.Get("/favicon.ico", func(c *echo.Context) error {
		c.Redirect(301, "/static/images/favicon.ico")
		return nil
	})

	// Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(
			app.Conf.UString("api.prefix"),
		),
	)

	// Bind React app
	app.Engine.Get("/*", app.React.Handle)
	return app
}

// Run runs the app
func (app *App) Run() {
	app.Engine.Run(":" + app.Conf.UString("port"))
}

// Custom renderer for Echo, to render html from bindata
type echoRenderer struct{}

func (er echoRenderer) Render(w io.Writer, name string, data interface{}) error {
	template := binhtml.New(Asset, AssetDir).MustLoadDirectory("templates")
	return template.Execute(w, data)
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
