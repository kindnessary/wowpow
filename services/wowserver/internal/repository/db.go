package repository

import (
	"context"
	"database/sql"

	"github.com/kindnessary/wowpow/services/wowserver/pkg/entity"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetQuote(ctx context.Context, id int) (*entity.Quote, error) {
	statement := `
		SELECT 
			"id",
			"text"
		FROM
			"quote"
		WHERE
			"id" = $1`

	var dbQuote quote
	err := r.db.QueryRowContext(ctx, statement, id).
		Scan(&dbQuote.id, &dbQuote.text)

	if err != nil {
		return nil, err
	}

	return &entity.Quote{
		ID:   dbQuote.id,
		Text: dbQuote.text,
	}, nil
}
