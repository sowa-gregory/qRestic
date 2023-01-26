package main

import (
	"fmt"
	"qrestic/gui"
	"qrestic/resticcmd"
)

func main() {
	gui := gui.NewGui()

	go func() {
		gui.ShowProgressInfinite()
		gui.DisableBackupButton()
		gui.SetStatus("reading snapshots")
		if data, err := resticcmd.GetSnapshots(); err == nil {
			gui.SetSnapshots(data)
			gui.SetStatus("OK")
			gui.EnableBackupButton()
		} else {
			gui.SetStatus("failed")
			gui.ShowError(err)
		}
		gui.HideProgress()

	}()

	gui.SetBackupCallback(func() {
		fmt.Println("he he")
		gui.ShowProgress(56)
	})

	gui.ShowAndRun()

}
