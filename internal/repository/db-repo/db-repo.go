package dbrepo

import (
	"database/sql"

	"github.com/parthvinchhi/bread-n-breakfast/internal/config"
	"github.com/parthvinchhi/bread-n-breakfast/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
