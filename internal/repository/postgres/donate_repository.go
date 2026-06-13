package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"pecatu-be/internal/domain"
)

type donateRepository struct {
	db *sql.DB
}

func NewDonateRepository(db *sql.DB) domain.DonateRepository {
	return &donateRepository{
		db: db,
	}
}

func (r *donateRepository) GetByDonateID(ctx context.Context, donateID string) ([]*domain.DonateDetail, error) {
	query := `
		SELECT 
			c_id, 
			c_donate_id, 
			c_cluster, 
			c_block, 
			c_number, 
			ts_created_at, 
			d_created_by, 
			d_amount, 
			c_name
		FROM 
			donate_detail
		WHERE 
			c_donate_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, donateID)
	if err != nil {
		return nil, fmt.Errorf("query donate_detail failed: %w", err)
	}
	defer rows.Close()

	var details []*domain.DonateDetail

	for rows.Next() {
		var detail domain.DonateDetail
		
		// Use sql.Null* types if any of these columns can be null in the database.
		// Assuming they are NOT NULL based on the initial requirements, 
		// but using basic types directly. If they can be null, this will return an error.
		err := rows.Scan(
			&detail.ID,
			&detail.DonateID,
			&detail.Cluster,
			&detail.Block,
			&detail.Number,
			&detail.CreatedAt,
			&detail.CreatedBy,
			&detail.Amount,
			&detail.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("scan donate_detail row failed: %w", err)
		}

		details = append(details, &detail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return details, nil
}
