package main

import (
	"fmt"
	"qrestic/gui"
	"qrestic/resticcmd"
)

func main2() {
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

func onProgress(data string) {
	fmt.Println("!" + data)
}

func main() {
	err := resticcmd.ExecuteCmdProgress("restic -r /tmp/rest backup /Users/sowisz/Programming/pythona /Users/sowisz/Programming/python --json",
		onProgress, "RESTIC_PASSWORD=aqq", "RESTIC_PROGRESS_FPS=2")
	if err != nil {
		fmt.Println(err)
	}

}
