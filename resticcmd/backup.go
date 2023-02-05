package resticcmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type BackupStatus struct {
	Message_type    string
	Seconds_elapsed int
	Percent_done    float64
	Total_files     int
	Files_done      int
	Total_bytes     int
	Bytes_done      int
}

type BackupSummary struct {
	Message_type          string
	Files_new             int
	Files_changed         int
	Files_unmodified      int
	Dirs_new              int
	Dirs_changed          int
	Dirs_unmodified       int
	Data_blobs            int
	Tree_blobs            int
	Data_added            int
	Total_files_processed int
	Total_bytes_processed int
	Total_duration        float64
	Snapshot_id           string
}

type StatusCallback func(status BackupStatus)
type SummaryCallback func(summary BackupSummary)

func doProgressCallbackDebug(line string, statusCb StatusCallback, summaryCb SummaryCallback) error {
	fmt.Println(line)
	return nil
}

func doProgressCallback(line string, statusCb StatusCallback, summaryCb SummaryCallback) error {
	if strings.HasPrefix(line, "{\"message_type\":\"status\"") {
		var status BackupStatus
		if err := json.Unmarshal([]byte(line), &status); err != nil {
			return err
		}
		statusCb(status)
	} else {
		var summary BackupSummary
		if err := json.Unmarshal([]byte(line), &summary); err != nil {
			return err
		}
		summaryCb(summary)
	}
	return nil
}

func processCmdStdOut(stdout io.ReadCloser, statusCb StatusCallback, summaryCb SummaryCallback) error {
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
				doProgressCallback(line.String(), statusCb, summaryCb)
				line.Reset()
			} else {
				line.WriteByte(buf[pos])
			}
			pos++
		}
	}
}

func executeCmdProgress(cmdLine string, statusCb StatusCallback, summaryCb SummaryCallback, env ...string) error {
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
	if err = processCmdStdOut(stdout, statusCb, summaryCb); err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("%w %s", err, errBuf.String())
	}
	return nil
}

func DoBackup(statusCb StatusCallback, summaryCb SummaryCallback) error {
	cmdLine := fmt.Sprintf("restic|-r|%s|backup|", configs[selectedConfig].Repository)
	for _, exc := range configs[selectedConfig].Excludes {
		cmdLine += "--exclude=" + exc + "|"
	}

	for _, src := range configs[selectedConfig].Sources {
		cmdLine += src + "|"
	}

	cmdLine += "--compression=max|--json"

	fmt.Println(cmdLine)
	err := executeCmdProgress(cmdLine, statusCb, summaryCb, "RESTIC_PASSWORD="+configs[selectedConfig].Password, "RESTIC_PROGRESS_FPS=2")
	if err != nil {
		return err
	}
	return nil
}
