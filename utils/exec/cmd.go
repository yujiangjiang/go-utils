package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Exec(cmdStr string, args ...string) (string, string) {

	if cmdStr == "" {
		return "", "empty cmd"
	}
	cmd := exec.Command(cmdStr, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Sprint("error when exec cmd: " + cmdStr + " " + strings.Join(args, " ") + ", " + stderr.String())
	}
	return strings.Trim(out.String(), " \n"), ""
}
