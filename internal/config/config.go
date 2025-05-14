package config

import (
	"os"
)

type Config struct {
	Port string
	DatabaseFile string
	TodoPassword string
	StaticPath string
}

func GetEnv() (*Config, error)  {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "../scheduler.db"
	}

	staticPathWeb := os.Getenv("STATIC_PATH")
	if staticPathWeb == "" {
		staticPathWeb = "web"
	}
	
	cfg := &Config{
		Port: port,
		DatabaseFile: dbFile,
		TodoPassword: "12345678",
		StaticPath: staticPathWeb,
	}

	return cfg, nil
}