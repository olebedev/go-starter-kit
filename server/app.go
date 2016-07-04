package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/itsjamie/go-bindata-templates"
	"github.com/nu7hatch/gouuid"
	"github.com/olebedev/config"
	"github.com/labstack/echo/engine/standard"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	Server *standard.Server
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

	// Set config variables delivered from main.go:11
	// Variables defined as ./conf.go:3
	conf.Set("debug", debug)
	conf.Set("commitHash", commitHash)

	// Parse environ variables for defined
	// in config constants
	conf.Env()

	// Make an engine
	e := echo.New()

	// Set up echo
	e.SetDebug(conf.UBool("debug"))

	// Regular middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/images/favicon.ico")
	})

	std := standard.WithConfig(engine.Config{})
	std.SetHandler(e)

	// Initialize the application
	app := &App{
		Conf:   conf,
		Engine: e,
		API:    &API{},
		React: NewReact(
			conf.UString("duktape.path"),
			conf.UBool("debug"),
			std,
		),
	}


	// Use precompiled embedded templates
	app.Engine.SetRenderer(NewTemplate())

	// Map app struct to access from request handlers
	// and middlewares
	app.Engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	// Map uuid for every requests
	app.Engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, _ := uuid.NewV4()
			c.Set("uuid", id)
			return next(c)
		}
	})

	////Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(
			app.Conf.UString("api.prefix"),
		),
	)

	// Create file http server from bindata
	fileServerHandler := http.FileServer(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
	})

	// Serve static via bindata and handle via react app
	// in case when static file was not found
	app.Engine.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// execute echo handlers chain
			err := next(c)
			// if page(handler) for url/method not found
			if err != nil && err.Error() == http.StatusText(http.StatusNotFound) {
				// check if file exists
				// omit first `/`
				if _, err := Asset(c.Request().URI()[1:]); err == nil {
					fileServerHandler.ServeHTTP(
						c.Response().(*standard.Response).ResponseWriter,
						c.Request().(*standard.Request).Request)
					return nil
				}
				// if static file not found handle request via react application
				return app.React.Handle(c)
			}
			// Move further if err is not `Not Found`
			return err
		}
	})

	return app
}

// Run runs the app
func (app *App) Run() {
	app.Engine.Run(standard.New(":" + app.Conf.UString("port")))
}

// Template is custom renderer for Echo, to render html from bindata
type Template struct {
	templates *template.Template
}

// NewTemplate creates a new template
func NewTemplate() *Template {
	return &Template{
		templates: binhtml.New(Asset, AssetDir).MustLoadDirectory("templates"),
	}
}

// Render renders template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context ) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// AppOptions is options struct
type AppOptions struct{}

func (ao *AppOptions) init() { /* write your own*/ }
