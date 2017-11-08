package category

import (
	"github.com/gregbiv/news-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

// NewStorer inits and returns an instance
// of category Storer
func NewStorer(db *sqlx.DB) Storer {
	return newCategoryManager(db, nil)
}

func (m *categoryManager) Store(d *model.Category) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	if err := m.insertCategory(tx, d); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (m *categoryManager) insertCategory(tx *sqlx.Tx, category *model.Category) error {
	query := `
        INSERT INTO category
        (
            category_id,
            name,
            title
        )
        VALUES ($1, $2, $3)
    `

	stmt, err := tx.Prepare(query)
	if err != nil {
		return stacktrace.Propagate(err, "failed to creates a prepared statement to store data in category table")
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		category.CategoryID,
		category.Name,
		category.Title,
	)

	return stacktrace.Propagate(err, "failed to store data into category table")
}
