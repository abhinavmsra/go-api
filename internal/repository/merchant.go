package repository

import "database/sql"

func (s *Storage) IndexMerchant() (*sql.Rows, error) {
	query := `SELECT * FROM "merchants" ORDER BY created_at DESC LIMIT 25;`
	return s.db.Query(query)
}

func (s *Storage) CreateMerchant(name string, secret string) *sql.Row {
	query := `INSERT INTO "merchants" (name, api_secret) VALUES ($1, $2) RETURNING *;`
	return s.db.QueryRow(query, name, secret)
}

func (s *Storage) ShowMerchant(id string) *sql.Row {
	query := `SELECT * FROM "merchants" WHERE "merchants"."id" = $1 LIMIT 1;`
	return s.db.QueryRow(query, id)
}

func (s *Storage) FindMerchantBySecret(secret string) *sql.Row {
	query := `SELECT * FROM "merchants" WHERE "merchants"."api_secret" = $1 LIMIT 1;`
	return s.db.QueryRow(query, secret)
}

func (s *Storage) UpdateMerchant(id string, name string) *sql.Row {
	query := `UPDATE "merchants" SET name = $2 WHERE "merchants"."id" = $1;`
	return s.db.QueryRow(query, id, name)
}

func (s *Storage) DeleteMerchant(id string) *sql.Row {
	query := `DELETE FROM "merchants" WHERE "merchants"."id" = $1 RETURNING *;`
	return s.db.QueryRow(query, id)
}
