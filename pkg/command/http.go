package command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/gregbiv/news-api/pkg/api/docs"
	"github.com/gregbiv/news-api/pkg/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// HTTPCommand is responsible for running the http server
type HTTPCommand struct {
	Meta
}

// Run is responsible for starting the http server
func (c *HTTPCommand) Run(args []string) int {
	flags := c.FlagSet("http")
	flags.Usage = func() { c.UI.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	router := chi.NewRouter()

	// Setup handler dependencies
	db, err := sqlx.Open("postgres", c.Config.Database.PostgresDB.DSN)
	if err != nil {
		log.Fatalf("Postgres Connection failed: %+v", err)
	}

	// Set package level configurations
	middleware.Debug = c.Config.Debug

	// A good base middleware stack
	router.Use(
		middleware.OpenTracing(opentracing.GlobalTracer()),
		chiMiddleware.WithValue(middleware.DatabaseConnection, db),
		chiMiddleware.WithValue("app.config", c.Config),
		chiMiddleware.Recoverer,
		lg.RequestLogger(logrus.StandardLogger()),
	)

	// HelloWorld
	router.Route("/", func(r chi.Router) {
		r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("News Today API"))
		})
	})

	// Documentation
	router.Route("/docs", docs.Docs)

	// Version 1
	router.Route("/v1", func(r chi.Router) {
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", c.Config.Port), router))

	return 0
}

// Help outputs a helper text for the command
func (*HTTPCommand) Help() string {
	helpText := `
Usage: news-api http [options]

  Start the Http Rest API server
`

	return strings.TrimSpace(helpText)
}

// Synopsis outputs to the console the synopsis of the command
func (c *HTTPCommand) Synopsis() string {
	return "Start the Http Rest API server"
}
