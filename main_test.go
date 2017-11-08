package main

import (
	"github.com/DATA-DOG/godog"
	"github.com/gregbiv/news-api/features/bootstrap"
	"github.com/gregbiv/news-api/pkg/config"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var (
	URL = "http://localhost"
)

func FeatureContext(s *godog.Suite) {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	bootstrap.RegisterGomega(s)

	conn, err := sqlx.Open("postgres", cfg.Database.PostgresDB.DSN)
	if err != nil {
		log.Fatal(err)
	}

	bootstrap.RegisterSystemContext(s, URL)

	bootstrap.RegisterCategoryContext(s, URL, conn)
}
