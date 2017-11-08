package bootstrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gregbiv/news-api/pkg/api"
	"github.com/jmoiron/sqlx"
	"github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
)

// RegisterCategoryContext Register the system context
func RegisterCategoryContext(s *godog.Suite, uri string, database *sqlx.DB) {
	category := &categoryContext{
		uri: uri,
		db:  database,
	}

	s.AfterScenario(category.resetResponse)

	s.Step(`^that "([^"]*)" category was created with title "([^"]*)":$`, category.thatCategoryWasCreatedWithTitle)
}

type categoryContext struct {
	uri      string
	response *http.Response
	db       *sqlx.DB
}

func (c *categoryContext) resetResponse(interface{}, error) {
	if c.response != nil {
		c.response.Body.Close()
		c.response = nil
	}
}

func assertErrorResponse(body io.ReadCloser, code string, target string, message string) error {
	actualBuff := new(bytes.Buffer)
	actualBuff.ReadFrom(body)
	actual := actualBuff.String()

	expectedError := api.ErrResponse{
		Errors: api.Error{
			Code:    code,
			Target:  target,
			Message: message,
		},
		HTTPStatusCode: http.StatusBadRequest,
	}

	expectedJSON, err := json.Marshal(expectedError)

	if !gomega.Expect(actual).Should(gomega.MatchJSON(expectedJSON)) {
		return errors.New("Invalid response")
	}

	return err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// GIVEN
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *categoryContext) thatCategoryWasCreatedWithTitle(name, title string, table *gherkin.DataTable) error {
	categoryID := uuid.NewV4().String()
	err := c.createCategory(categoryID, name, title)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HELPERS
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type categoryRequest struct {
	Response string            `json:"response"`
	Result   []categoryContext `json:"result"`
}

func (c *categoryContext) CreateCategory(category categoryContext) (err error) {
	data, err := json.Marshal(category)
	if err != nil {
		return
	}

	contentReader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", c.uri+"/category", contentReader)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	c.response, err = client.Do(req)

	return
}

func (c *categoryContext) createCategory(categoryID string, name string, title string) error {
	// Inser category
	_, err := c.db.Exec(
		`INSERT INTO category
		(
			category_id,
			name,
			title
		)
		VALUES ($1, $2, $3)`,
		categoryID,
		name,
		title,
	)

	if err != nil {
		return err
	}

	return nil
}
