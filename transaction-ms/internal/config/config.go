package config

import (
	"time"

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
	Database
	Grpc
}

type Grpc struct {
	GrpcHost string `env:"GRPC_HOST"`
	GrpcPort string `env:"GRPC_PORT"`
}

type Database struct {
	DbUser     string        `env:"DB_USER" envDefault:""`
	DbPassword string        `env:"DB_PASSWORD" envDefault:""`
	DbHost     string        `env:"DB_HOST" envDefault:""`
	DbName     string        `env:"DB_NAME" envDefault:""`
	DbOptions  string        `env:"DB_OPTIONS" envDefault:""`
	Timeout    time.Duration `env:"DB_TIMEOUT" envDefault:"10s"`
}
