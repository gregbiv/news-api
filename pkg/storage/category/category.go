package category

import (
	"errors"
	"github.com/gregbiv/news-api/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
)

var (
	// ErrCategoryNotFound ...
	ErrCategoryNotFound = errors.New("Unknown category")
)

type (
	// Storer describes logic
	// for persisting a category
	Storer interface {
		Store(category *model.Category) error
	}

	// Discarder describes logic
	// for discarding a category
	Discarder interface {
		Discard(ID string) error
	}

	// Getter is the object responsible for getting a category
	Getter interface {
		// GetCategoryBySubscriptionAndInterval gets a category model from the database given an subscriptionId and DeliveryIntervalId
		GetCategoryByID(ID string) (*model.Category, error)
	}

	// Updater is the object responsible for updating a category
	Updater interface {
		//Update updates a category in the database given an updated category model
		Update(category *model.Category) error
	}

	// Asserter is the object responsible for asserting a category
	Asserter interface {
		AssertExists(ID string) (bool, error)
	}

	// categoryManager handlers
	// category Store, Publish and Discard
	// operations
	categoryManager struct {
		db               *sqlx.DB
		categoryAsserter Asserter
	}

	// Category describes a category db model
	category struct {
		CategoryID uuid.UUID `db:"category_id"`
		Name       string    `db:"name"`
		Title      string    `db:"title"`
	}
)

// newCategoryManager inits and returns
// an instance of category manager
func newCategoryManager(
	db *sqlx.DB,
	categoryAsserter Asserter,
) *categoryManager {
	return &categoryManager{
		db:               db,
		categoryAsserter: categoryAsserter,
	}
}
