package cli

import (
	"fmt"
	"os/exec"
)

func CheckCommandExists(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("can't find '%s' in the PATH. Install it first", command)
	}

	return nil
}
