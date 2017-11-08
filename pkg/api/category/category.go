package category

import (
	"encoding/json"
	"errors"
	"github.com/gregbiv/news-api/pkg/model"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
)

// ErrInvalidBody represents the error when the request body is invalid.
var ErrInvalidBody = errors.New("Invalid request body provided")

type (
	// category describes a news-api API model
	category struct {
		CategoryID *uuid.UUID `json:"category_id"`
		Name       string     `json:"name"`
		Title      string     `json:"title"`
	}
)

func (c *category) fromDB(dbCategory *model.Category) error {
	id, err := uuid.FromString(dbCategory.CategoryID)
	if err != nil {
		return err
	}

	c.CategoryID = &id
	c.Name = dbCategory.Name
	c.Title = dbCategory.Title

	return nil
}

func (c *category) toModel() (modelCategory model.Category) {
	modelCategory.CategoryID = c.CategoryID.String()
	modelCategory.Name = c.Name
	modelCategory.Title = c.Title

	return modelCategory
}

func (c *category) fromRequest(r *http.Request) error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ErrInvalidBody
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}

	return nil
}
