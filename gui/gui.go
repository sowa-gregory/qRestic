package gui

import (
	"qrestic/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Gui struct {
	app               fyne.App
	treeData          types.SnapshotTree
	window            fyne.Window
	tree              *widget.Tree
	progressInfinite  *widget.ProgressBarInfinite
	progress          *widget.ProgressBar
	status            *widget.Label
	backupButton      *widget.Button
	combo             *widget.Select
	comboPrevSelected int

	backupCallback *func()
	comboCallback  *func(index int)
}

func NewGui() *Gui {
	var gui Gui
	gui.createGui()
	return &gui
}

func (gui *Gui) SetBackupCallback(callback func()) {
	gui.backupCallback = &callback
}

func (gui *Gui) SetSnapshots(data types.SnapshotTree) {
	gui.treeData = data
	gui.tree.Refresh()
}

func (gui *Gui) EnableBackupButton() {
	gui.backupButton.Enable()
}

func (gui *Gui) DisableBackupButton() {
	gui.backupButton.Disable()
}

func (gui *Gui) EnableCombo() {
	gui.combo.Enable()
}

func (gui *Gui) DisableCombo() {
	gui.combo.Disable()
}

func (gui *Gui) ShowProgress(value float64) {
	gui.progress.SetValue(value)
	gui.progress.Show()
	gui.progressInfinite.Hide()
}

func (gui *Gui) SetStatus(status string) {
	gui.status.SetText("Status: " + status)
}

func (gui *Gui) ShowProgressInfinite() {

	gui.progress.Hide()
	gui.progressInfinite.Show()
}

func (gui *Gui) HideProgress() {
	gui.progress.Hide()
	gui.progressInfinite.Hide()
}

func (gui *Gui) ShowAndRun() {
	gui.window.ShowAndRun()
}

func (gui *Gui) ShowError(err error, quit bool) {
	d := dialog.NewError(err, gui.window)
	if quit {
		d.SetOnClosed(gui.app.Quit)
	}
	d.Show()
}

func (gui *Gui) SetCombo(options []string) {
	gui.combo.Options = options
	gui.combo.SetSelectedIndex(0)
	gui.comboPrevSelected = 0
	gui.combo.Refresh()
}

func (gui *Gui) SetComboCallback(callback func(index int)) {
	gui.comboCallback = &callback
}

func (gui *Gui) createGui() {
	gui.app = app.New()

	gui.window = gui.app.NewWindow("qRestic v" + AppVersion)
	gui.treeData = make(types.SnapshotTree)

	gui.tree = &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return gui.treeData[uid]
		},
		IsBranch: func(uid string) bool {
			_, b := gui.treeData[uid]
			return b
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Template Object")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			node.(*widget.Label).SetText(uid)
		},
	}

	gui.progress = widget.NewProgressBar()
	gui.progress.Hide()
	gui.progressInfinite = widget.NewProgressBarInfinite()
	gui.progressInfinite.Hide()
	gui.status = widget.NewLabel("")
	gui.backupButton = widget.NewButton("Backup", func() {
		if gui.backupCallback != nil {
			(*gui.backupCallback)()
		}
	})

	gui.combo = widget.NewSelect(nil, func(value string) {
		sel := gui.combo.SelectedIndex()
		if gui.comboCallback != nil && sel != gui.comboPrevSelected {
			gui.comboPrevSelected = sel
			(*gui.comboCallback)(sel)
		}
	})
	gui.combo.SetSelectedIndex(0)

	topBox := container.NewGridWithColumns(3, widget.NewLabel("Config:"), gui.combo, container.New(layout.NewMaxLayout(), gui.progressInfinite, gui.progress))
	content := container.NewBorder(topBox,
		container.NewVBox(widget.NewSeparator(), gui.backupButton, gui.status),
		nil, nil, widget.NewCard("Snapshots", "", gui.tree))

	//content := container.NewBorder(nil, nil, nil, buttonContent, listContent)

	gui.window.SetContent(content)
	gui.window.Resize(fyne.NewSize(800, 600))
}
