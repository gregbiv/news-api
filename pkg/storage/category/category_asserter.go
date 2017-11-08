package category

import (
	"github.com/jmoiron/sqlx"
)

// dbCategoryAsserter implements Asserter interface
type dbCategoryAsserter struct {
	db sqlx.Queryer
}

// NewAsserter inits and returns an instance of category asserter
func NewAsserter(db sqlx.Queryer) Asserter {
	return &dbCategoryAsserter{db}
}

func (s *dbCategoryAsserter) AssertExists(ID string) (bool, error) {
	query := `
		SELECT
			count(category_id) as total
		FROM
			category
		WHERE
			category_id = $1`

	var total int
	err := s.db.QueryRowx(
		query,
		ID,
	).Scan(&total)

	if err != nil || total == 0 {
		return false, err
	}

	return true, nil
}
