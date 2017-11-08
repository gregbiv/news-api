package category

import (
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
)

// NewDiscarder inits and returns
// an instance of category Discarder
func NewDiscarder(db *sqlx.DB) Discarder {
	return newCategoryManager(db, nil)
}

func (m *categoryManager) Discard(ID string) error {
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	if err := discardCategory(tx, ID); err != nil {
		if err == ErrCategoryNotFound {
			return err
		}
		log.Panic(stacktrace.Propagate(err, "Failed to discard category", ID))
	}

	if err := tx.Commit(); err != nil {
		log.Panic(stacktrace.Propagate(err, "Failed to commit transaction to delete a category (Category ID: %s)", ID))
	}

	return nil
}

func discardCategory(tx *sqlx.Tx, ID string) error {
	query := `
		DELETE
			FROM category
			WHERE category_id = $1`

	result, err := tx.Exec(query, ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
