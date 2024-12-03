package main

import (
	"server/internal/config"
	"server/internal/repository/connector"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.Load(); err != nil {
		logrus.Fatal("couldn't load env configs", err)
		return
	}
	dbConfig, err := config.NewDBConfig()
	if err != nil {
		logrus.Fatal("couldn't load db config", err)
		return
	}
	logrus.Info(dbConfig)
	_, err = connector.NewDatabase(dbConfig)
	if err != nil {
		logrus.Fatal("couldn't make db connection", err)
		return
	}

}
