package server

import (
	"io"
	"net/http"
	"regexp"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/itsjamie/go-bindata-templates"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
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

	// Middlewares
	engine.Use(mw.Logger())
	engine.Use(mw.Recover())

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

	// Load embedded templates MISSING
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

	var staticPathRegExp = regexp.MustCompile(`^/static/?.*`)

	staticAsset := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
		Prefix:   "",
	}

	fileServerHandler := http.FileServer(&staticAsset)

	// Handle all not found routes via react app (except for static files)
	app.Engine.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := h(c)
			if err != nil && err.Error() == http.StatusText(http.StatusNotFound) {
				if staticPathRegExp.MatchString(c.Request().URL.Path) {
					// Handles static files with assetFs
					fileServerHandler.ServeHTTP(c.Response(), c.Request())
					return nil
				}
				//if not found and not a static file then serv with react.
				return app.React.Handle(c)
			}
			return err
		}
	})

	// Avoid favicon react handling
	app.Engine.Get("/favicon.ico", func(c *echo.Context) error {
		c.Redirect(301, "/static/images/favicon.ico")
		return nil
	})

	//Force handle index to React otherwise it will throw Method not allowed.
	app.Engine.Get("/", app.React.Handle)

	// Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(app.Conf.UString("api.prefix")),
	)
	return app
}

// Run runs the app
func (app *App) Run() {
	app.Engine.Run(":" + app.Conf.UString("port"))
}

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
