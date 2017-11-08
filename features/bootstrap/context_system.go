package bootstrap

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/onsi/gomega"
)

// RegisterSystemContext Register the system context
func RegisterSystemContext(s *godog.Suite, uri string) {
	sys := &systemContext{uri: strings.TrimRight(uri, "/")}

	s.AfterScenario(sys.resetResponse)

	s.Step(`^I send a "(GET|POST|PUT|DELETE)" request to "([^"]*)"$`, sys.iSendRequestTo)
	s.Step(`^I send a "(POST|PUT)" request to "([^"]*)" containing the following JSON:$`, sys.iSendARequestToContainingTheFollowingJSON)
	s.Step(`^the response code should be (\d+)$`, sys.theResponseCodeShouldBe)
	s.Step(`^the response should match json:$`, sys.theResponseShouldMatchJSON)
	s.Step(`^the response should match csv:$`, sys.theResponseShouldMatchCsv)

	s.Step(`^I request the status url of the system$`, sys.iRequestTheStatusURLOfTheSystem)
	s.Step(`^I should receive a response with all status checks successful$`, sys.iShouldReceiveAResponseWithAllStatusChecksSuccessful)

	s.Step(`^I request a url that does not exist in the system$`, sys.iRequestAUrlThatDoesNotExistInTheSystem)
	s.Step(`^I should see that the request page cannot be found$`, sys.iShouldSeeThatTheRequestPageCannotBeFound)

	s.Step(`^I request the api documentation$`, sys.iRequestTheAPIDocumentation)
	s.Step(`^I should see that the api documentation$`, sys.iShouldSeeThatTheAPIDocumentation)
}

type systemContext struct {
	uri      string
	response *http.Response
}

func (c *systemContext) resetResponse(interface{}, error) {
	if c.response != nil {
		c.response.Body.Close()
		c.response = nil
	}
}

func (c *systemContext) iSendRequestTo(method, endpoint string) (err error) {
	req, err := http.NewRequest(method, c.uri+endpoint, nil)
	if err != nil {
		return
	}

	c.response, err = http.DefaultClient.Do(req)
	return
}

func (c *systemContext) iSendARequestToContainingTheFollowingJSON(
	method string,
	endpoint string,
	JSON *gherkin.DocString,
) (err error) {
	req, err := http.NewRequest(method, c.uri+endpoint, strings.NewReader(JSON.Content))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	c.response, err = http.DefaultClient.Do(req)
	return
}

func (c *systemContext) theResponseCodeShouldBe(code int) error {
	if code != c.response.StatusCode {
		body, _ := ioutil.ReadAll(c.response.Body)
		return fmt.Errorf(
			"expected response code to be: %d, but actual is: %d with body %s",
			code,
			c.response.StatusCode,
			body,
		)
	}
	return nil
}

func (c *systemContext) theResponseShouldMatchJSON(body *gherkin.DocString) error {
	contentType := c.response.Header.Get("Content-type")
	if contentType != "application/json" {
		return fmt.Errorf("expected Content-Type to be: %s, but actual is: %s", "application/json", contentType)
	}

	actualBuff := new(bytes.Buffer)
	actualBuff.ReadFrom(c.response.Body)
	actual := actualBuff.String()

	if !gomega.Expect(actual).Should(gomega.MatchJSON(body.Content)) {
		return errors.New("Invalid response")
	}

	return nil
}

func (c *systemContext) theResponseShouldMatchCsv(body *gherkin.DocString) error {
	contentType := c.response.Header.Get("Content-type")
	if contentType != "text/csv" {
		return fmt.Errorf("expected Content-Type to be: %s, but actual is: %s", "text/csv", contentType)
	}

	actualBuff := new(bytes.Buffer)
	actualBuff.ReadFrom(c.response.Body)
	actual := strings.TrimRight(actualBuff.String(), "\n")

	expected := strings.TrimRight(body.Content, "\n")

	if !gomega.Expect(actual).Should(gomega.Equal(expected)) {
		return errors.New("Invalid response")
	}

	return nil
}

func (c *systemContext) iRequestTheStatusURLOfTheSystem() error {
	return c.iSendRequestTo("GET", "/status")
}

func (c *systemContext) iShouldReceiveAResponseWithAllStatusChecksSuccessful() error {
	return c.theResponseCodeShouldBe(http.StatusOK)
}

func (c *systemContext) iRequestAUrlThatDoesNotExistInTheSystem() error {
	return c.iSendRequestTo("GET", "/UrlThatDoesNotExistInTheSystem")
}

func (c *systemContext) iShouldSeeThatTheRequestPageCannotBeFound() error {
	return c.theResponseShouldMatchJSON(&gherkin.DocString{
		Content: `
        {
            "error": {
                "code":"InvalidUri",
                "message":"The requested URI does not represent any resource on the server."
            }
        }`,
	})
}
func (c *systemContext) iRequestTheAPIDocumentation() error {
	return c.iSendRequestTo("GET", "/docs/api.raml")
}

func (c *systemContext) iShouldSeeThatTheAPIDocumentation() error {
	return c.theResponseCodeShouldBe(http.StatusOK)
}
