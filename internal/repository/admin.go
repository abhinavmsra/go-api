package repository

import "database/sql"

func (s *Storage) FindAdminBySecret(secret string) *sql.Row {
	query := `SELECT * FROM "admins" WHERE "admins"."api_secret" = $1 LIMIT 1;`
	return s.db.QueryRow(query, secret)
}

func (s *Storage) FindAdminByName(name string) *sql.Row {
	query := `SELECT * FROM "admins" WHERE "admins"."name" = $1 LIMIT 1;`
	return s.db.QueryRow(query, name)
}

func (s *Storage) CreateAdmin(name string, secret string) *sql.Row {
	query := `INSERT INTO "admins" (name, api_secret) VALUES ($1, $2) RETURNING *;`
	return s.db.QueryRow(query, name, secret)
}
