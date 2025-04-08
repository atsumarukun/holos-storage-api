package api

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDatabase(conf *databaseConfig) (*sqlx.DB, error) {
	c := &mysql.Config{
		Addr:      conf.Host + ":" + conf.Port,
		User:      conf.User,
		Passwd:    conf.Password,
		DBName:    conf.Database,
		Net:       "tcp",
		ParseTime: true,
	}

	db, err := sqlx.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
