package category

import (
	"database/sql"
	"github.com/gregbiv/news-api/pkg/model"
	"github.com/jmoiron/sqlx"
)

type (
	dbGetter struct {
		db *sqlx.DB
	}
)

// NewGetter inits and returns a Getter instance
func NewGetter(db *sqlx.DB) Getter {
	return &dbGetter{db: db}
}

func (dg *dbGetter) GetCategoryByID(ID string) (*model.Category, error) {
	dbCategory := category{}

	query := `
        SELECT
            category_id,
			name,
			title
        FROM
            category
		WHERE
			category_id = $1
        ORDER BY order DESC
	`

	err := dg.db.Get(&dbCategory, query, ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	modelCategory := &model.Category{}
	modelCategory.CategoryID = dbCategory.CategoryID.String()
	modelCategory.Name = dbCategory.Name
	modelCategory.Title = dbCategory.Title

	return modelCategory, nil
}
