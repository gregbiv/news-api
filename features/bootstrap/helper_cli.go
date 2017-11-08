package bootstrap

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func commandExecute(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stdErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return fmt.Errorf("Failed to read the output: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Failed to execute the command (%s): %s", err, stdErr)
	}

	return nil
}
