package main

import (
	"os"

	"github.com/gregbiv/news-api/pkg/command"
	"github.com/gregbiv/news-api/pkg/command/migration"
	"github.com/gregbiv/news-api/pkg/config"
	"github.com/mattes/migrate/database"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{FieldMap: log.FieldMap{
		log.FieldKeyTime: "@timestamp",
		log.FieldKeyMsg:  "message",
	}})

	// Config
	conf, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Logging
	level, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(level)
	if conf.Debug {
		log.Debugf("Initialized with config: %+v", conf)
	}

	db, err := database.Open(conf.Database.PostgresDB.DSN)
	if err != nil {
		log.Fatalf("Cannot initialize db with the following DSN: %s", conf.Database.PostgresDB.DSN)
	}

	c := &cli.CLI{
		Name:     "news-api",
		Version:  "dev",
		HelpFunc: cli.BasicHelpFunc("news-api"),
		Commands: commands(conf, db),
		Args:     os.Args[1:],
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func commands(conf *config.Specification, db database.Driver) map[string]cli.CommandFactory {
	meta := command.Meta{
		UI: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Config: conf,
	}

	cf := map[string]cli.CommandFactory{
		"http": func() (cli.Command, error) {
			return &command.HTTPCommand{
				Meta: meta,
			}, nil
		},
		"migrate": func() (cli.Command, error) {
			cmd := migration.CreateMigrateCommand(
				db,
				conf.Migration.Dir,
				conf.Migration.Version,
			)
			return cmd, nil
		},
	}

	return cf
}
