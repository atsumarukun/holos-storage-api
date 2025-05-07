package api

import "os"

type serverConfig struct {
	database   databaseConfig
	fileSystem fileSystemConfig
}

func loadServerConfig() *serverConfig {
	return &serverConfig{
		database:   *loadDatabaseConfig(),
		fileSystem: *loadFileSystemConfig(),
	}
}

type databaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func loadDatabaseConfig() *databaseConfig {
	return &databaseConfig{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Database: os.Getenv("DATABASE_NAME"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
	}
}

type fileSystemConfig struct {
	BasePath string
}

func loadFileSystemConfig() *fileSystemConfig {
	return &fileSystemConfig{
		BasePath: os.Getenv("FILE_SYSTEM_BASE_PATH"),
	}
}
