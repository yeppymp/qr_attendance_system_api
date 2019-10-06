package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"main.go/api"
	"main.go/app"
	"net/http"
	"strconv"
	"time"
)

var configFile string

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Warnf("Unable to read config from file")
	}
}

func main() {
	initConfig()

	a, appErr := app.New()
	if appErr != nil {
		logrus.Error(appErr)
	}
	defer a.Close()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter()
	api.Init(router.PathPrefix("/api").Subrouter())

	s := &http.Server{
		Addr:           ":" + strconv.Itoa(viper.GetInt("server.port")),
		Handler:        cors(router),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logrus.Infof("Serving api at http://127.0.0.1:%d", viper.GetInt("server.port"))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Error(err)
	}
}
