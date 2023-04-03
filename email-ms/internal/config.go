package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Config config

func init() {
	// Loading the environment variables from '.env' file.
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Error you can load the env variables by file")
	}

	if err := env.Parse(&Config); err != nil {
		logrus.Fatalf("Error initializing: %s", err.Error())
	}
}

type config struct {
	Environment string `env:"APP_ENV"`
	AppPort     string `env:"APP_PORT"`
	Grpc
	Mail
}

type Grpc struct {
	GrpcHost string `env:"GRPC_HOST"`
	GrpcPort string `env:"GRPC_PORT"`
}

type Mail struct {
	EmailHost        string `env:"EMAIL_HOST"`
	EmailPort        int    `env:"EMAIL_PORT"`
	EmailDomain      string `env:"EMAIL_DOMAIN"`
	EmailPassword    string `env:"EMAIL_PASSWORD"`
	EmailUserName    string `env:"EMAIL_USERNAME"`
	EmailEncription  string `env:"EMAIL_ENCRIPTION"`
	EmailFromName    string `env:"EMAIL_FROM_NAME"`
	EmailFromAddress string `env:"EMAIL_FROM_ADDRESS"`
}
