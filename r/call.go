package r

import (
	"os/exec"
)

func Call(cmd string) ([]byte, error) {
	out, err := exec.Command(
		"R",
		"-s",
		"-e",
		cmd,
	).Output()

	return out, err
}

func IsPackage(pak string) bool {
	output, err := Call(
		"x <- requireNamespace('" + pak + "');cat(tolower(x));",
	)

	if err != nil {
		return false
	}

	return string(output) == "true"
}
