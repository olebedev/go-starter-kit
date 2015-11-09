package server

import (
	"io"
	"net/http"

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
	if !conf.UBool("debug") {
		engine.SetDebug(true)
	}

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

	// Define static dir
	app.Engine.Static("/static", "/client")

	// Load embedded templates MISSING
	app.Engine.SetRenderer(echoRenderer{})

	// Map app struct to access from request handlers
	// and middlewares
	var contextSetterMiddleware func(c *echo.Context) error
	contextSetterMiddleware = func(c *echo.Context) error {
		c.Set("app", app)
		return nil
	}

	app.Engine.Use(contextSetterMiddleware)

	var requestIDMiddleware func(c *echo.Context) error
	requestIDMiddleware = func(c *echo.Context) error {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
		return nil
	}

	// Map uuid for every requests
	app.Engine.Use(requestIDMiddleware)

	// Avoid favicon react handling
	app.Engine.Get("/favicon.ico", func(c *echo.Context) error {
		c.Redirect(301, "/static/images/favicon.ico")
		return nil
	})

	// Bind api hadling for URL api.prefix
	app.API.Bind(
		app.Engine.Group(app.Conf.UString("api.prefix")),
	)

	// Handle all not found routes via react app
	app.Engine.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		if err.Error() == http.StatusText(http.StatusNotFound) {
			app.React.Handle(c)
		}
	})

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
