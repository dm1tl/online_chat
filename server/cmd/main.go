package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"server/clients/sso"
	"server/internal/config"
	"server/internal/handler"
	"server/internal/hub"
	"server/internal/repository"
	"server/internal/repository/connector"
	"server/internal/services"
	"syscall"

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
	db, err := connector.NewDatabase(dbConfig)
	if err != nil {
		logrus.Fatal("couldn't make db connection", err)
		return
	}
	ssocfg, err := config.NewSSOConfig()
	if err != nil {
		logrus.Fatal("couldn't load sso_config", err)
		return
	}
	ssogrpcServiceClient, err := sso.NewSSOServiceClient(logrus.New(), *ssocfg)
	if err != nil {
		logrus.Fatal("couldn't create ssoServiceClient", err)
		return
	}
	ssoClient := sso.NewSSOClientWrapper(ssogrpcServiceClient)
	repository := repository.NewRepository(db.GetDB())
	service := services.NewService(repository, ssoClient)
	wshub := hub.NewHub()
	go wshub.Run()
	wsHandler := hub.NewWSHandler(wshub)
	handler := handler.NewHandler(service, wsHandler)
	httpServerCfg, err := config.NewHTTPServerConfig()
	if err != nil {
		logrus.Fatal("couldn't get http server config", err)
		return
	}
	server := config.NewServer(httpServerCfg, handler.InitRoutes())
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("cannot start server %s", err.Error())
		}
	}()
	logrus.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("online chat shutting down")
	if err := server.ShutDown(context.Background()); err != nil {
		logrus.Errorf("couldn't shut down an online chat %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("couldn't close db connection %s", err.Error())
	}
	logrus.Print("online chat shutted down")

}
