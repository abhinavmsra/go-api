package repository

import "database/sql"

func (s *Storage) FindAdminBySecret(secret string) *sql.Row {
	query := `SELECT * FROM "admins" WHERE "admins"."api_secret" = $1 LIMIT 1;`
	return s.db.QueryRow(query, secret)
}
