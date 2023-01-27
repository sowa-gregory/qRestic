package main

import (
	"fmt"
	"qrestic/gui"
	"qrestic/resticcmd"
)

var g *gui.Gui

func onBackupStatus(status resticcmd.BackupStatus) {
	g.SetStatus(fmt.Sprintf("storing %d/%d files, elapsed:%d s", status.Files_done, status.Total_files, status.Seconds_elapsed))
	g.ShowProgress(status.Percent_done)
}

func onBackupSummary(summary resticcmd.BackupSummary) {
	readingSnaphots()
	g.SetStatus(fmt.Sprintf("Done - %d new files, %d changed files, elapsed:%d s", summary.Files_new, summary.Files_changed, (int)(summary.Total_duration)))
	g.HideProgress()

}

func readingSnaphots() {
	g.ShowProgressInfinite()
	g.DisableBackupButton()
	g.SetStatus("reading snapshots")
	if data, err := resticcmd.GetSnapshots(); err == nil {
		g.SetSnapshots(data)
		g.SetStatus("OK")
		g.EnableBackupButton()
	} else {
		g.SetStatus("failed")
		g.ShowError(err)
	}
	g.HideProgress()
}

func main() {
	g = gui.NewGui()

	go readingSnaphots()

	g.SetBackupCallback(func() {
		g.DisableBackupButton()

		resticcmd.DoBackup(onBackupStatus, onBackupSummary)

	})

	g.ShowAndRun()

}
