package api

import "os"

type serverConfig struct {
	database   databaseConfig
	auth       authConfig
	fileSystem fileSystemConfig
}

func loadServerConfig() *serverConfig {
	return &serverConfig{
		database:   loadDatabaseConfig(),
		auth:       loadAuthConfig(),
		fileSystem: loadFileSystemConfig(),
	}
}

type databaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func loadDatabaseConfig() databaseConfig {
	return databaseConfig{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Database: os.Getenv("DATABASE_NAME"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
	}
}

type authConfig struct {
	Endpoint string
}

func loadAuthConfig() authConfig {
	return authConfig{
		Endpoint: os.Getenv("AUTH_API_ENDPOINT"),
	}
}

type fileSystemConfig struct {
	BasePath string
}

func loadFileSystemConfig() fileSystemConfig {
	return fileSystemConfig{
		BasePath: os.Getenv("FILE_SYSTEM_BASE_PATH"),
	}
}
