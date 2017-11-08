package bootstrap

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func assertStatusEquals(response *http.Response, statusCode int) error {
	if response.StatusCode == statusCode {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Invalid response. Error when trying to retrieve the body. StatusCode: %d ", response.StatusCode)
	}

	return fmt.Errorf("Invalid response. StatusCode: %d Body: %s", response.StatusCode, body)
}

func assertNotFoundResponse(response *http.Response) error {
	if err := assertStatusEquals(response, http.StatusNotFound); err != nil {
		return err
	}
	return assertErrorResponse(response.Body, "InvalidUri", "", "The requested URI does not represent any resource on the server.")
}
