package bootstrap

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/onsi/gomega"
)

// RegisterGomega Enable gomega validation support
func RegisterGomega(s *godog.Suite) {
	failures := []string{}

	gomega.RegisterFailHandler(func(message string, callerSkip ...int) {
		failures = append(failures, message)
	})

	s.AfterStep(func(step *gherkin.Step, err error) {
		if err == nil {
			return
		}

		for _, failure := range failures {
			fmt.Printf("%s\n", failure)
		}
		failures = []string{}
	})
}
