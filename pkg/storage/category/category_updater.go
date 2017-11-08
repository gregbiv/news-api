package category

import (
	"errors"
	"github.com/gregbiv/news-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/palantir/stacktrace"
)

//ErrUpdatingCategory is being thrown when something goes wrong while updating the category
var ErrUpdatingCategory = errors.New("Something went wrong when trying to update the category")

type (
	dbCategoryUpdater struct {
		db *sqlx.DB
	}
)

// NewUpdater inits and returns a CategoryUpdater instance
func NewUpdater(db *sqlx.DB) Updater {
	return &dbCategoryUpdater{db: db}
}

func (du *dbCategoryUpdater) Update(category *model.Category) error {
	tx, err := du.db.Beginx()
	if err != nil {
		return err
	}

	steps := []func(*sqlx.Tx, *model.Category) (bool, error){
		du.updateCategory,
	}

	for _, step := range steps {
		ok, err := step(tx, category)
		if err != nil {
			tx.Rollback()
			stacktrace.Propagate(err, ErrUpdatingCategory.Error())
			return err
		}

		if !ok {
			tx.Rollback()
			return ErrUpdatingCategory
		}
	}

	return tx.Commit()
}

func (du *dbCategoryUpdater) updateCategory(tx *sqlx.Tx, category *model.Category) (bool, error) {
	query := `
        UPDATE
			category
		SET
			(
				name,
				title
			)
        =
			($1, $2)
		WHERE
			category_id = $3
    `

	return du.executeQuery(
		tx,
		query,
		category.Name,
		category.Title,
		category.CategoryID,
	)
}

func (du *dbCategoryUpdater) executeQuery(tx *sqlx.Tx, query string, args ...interface{}) (bool, error) {
	result, err := tx.Exec(query, args...)
	if err != nil {
		return false, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
