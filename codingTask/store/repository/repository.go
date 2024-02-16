package repository

import "database/sql"

type RepoPostgres struct {
	db *sql.DB
}
