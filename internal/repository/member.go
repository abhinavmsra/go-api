package repository

import "database/sql"

func (s *Storage) IndexMember(merchantId int, limit int, offset int) (*sql.Rows, error) {
	query := `SELECT * FROM "members" WHERE "members"."merchant_id" = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;`
	return s.db.Query(query, merchantId, limit, offset)
}

func (s *Storage) CreateMember(name string, secret string, email string, merchantId int) *sql.Row {
	query := `INSERT INTO "members" (name, api_secret, email, merchant_id) VALUES ($1, $2, $3, $4) RETURNING *;`
	return s.db.QueryRow(query, name, secret, email, merchantId)
}

func (s *Storage) ShowMember(id string, merchantId int) *sql.Row {
	query := `SELECT * FROM "members" WHERE "members"."id" = $1 AND "members"."merchant_id" = $2 LIMIT 1;`
	return s.db.QueryRow(query, id, merchantId)
}

func (s *Storage) UpdateMember(id string, merchantId int, name string) *sql.Row {
	query := `UPDATE "members" SET name = $3 WHERE "members"."id" = $1 AND "members"."merchant_id" = $2;`
	return s.db.QueryRow(query, id, merchantId, name)
}

func (s *Storage) DeleteMember(id string, merchantId int) *sql.Row {
	query := `DELETE FROM "members" WHERE "members"."id" = $1 AND "members"."merchant_id" = $2 RETURNING *;`
	return s.db.QueryRow(query, id, merchantId)
}
