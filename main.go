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
		g.ShowError(err, true)
	}
	g.HideProgress()
}

func main() {
	g = gui.NewGui()
	resticcmd.ReadConfig("backup-conf.json")

	go readingSnaphots()

	g.SetBackupCallback(func() {
		g.DisableBackupButton()

		if err := resticcmd.DoBackup(onBackupStatus, onBackupSummary); err != nil {
			g.SetStatus("failed")
			g.ShowError(err, false)
		}
		g.EnableBackupButton()
	})

	g.ShowAndRun()

}
