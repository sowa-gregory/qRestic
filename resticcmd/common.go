package resticcmd

import (
	"os"
	"os/exec"
	"strings"
)

func prepareCmd(cmdLine string, env ...string) *exec.Cmd {
	cmdArgs := strings.Split(cmdLine, "|")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = append(os.Environ(), env...)
	return cmd
}
