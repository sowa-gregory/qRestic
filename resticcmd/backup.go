package resticcmd

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func processCmdStdOut(stdout io.ReadCloser, progressCallback func(progress string)) error {
	buf := make([]byte, 1024)
	var line strings.Builder

	for {
		readBytes, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		pos := 0
		for pos < readBytes {
			if buf[pos] == '\n' && line.Len() > 0 {
				progressCallback(line.String())
				line.Reset()
			} else {
				line.WriteByte(buf[pos])
			}
			pos++
		}
	}
}

func ExecuteCmdProgress(cmdLine string, progressCallback func(progress string), env ...string) error {
	cmd := prepareCmd(cmdLine, env...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		return err
	}
	if err = processCmdStdOut(stdout, progressCallback); err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("%w %s", err, errBuf.String())
	}
	return nil
}
