package main

import (
	"fmt"
	"os/exec"

	"github.com/rivo/tview"
)

var (
	preview *tview.TextView

	displayPrim = tview.NewDropDown().SetLabel("Display:")
	onPrim = tview.NewCheckbox().SetLabel("On:")
	modePrim = tview.NewDropDown().SetLabel("Mode:")
	primaryPrim = tview.NewCheckbox().SetLabel("Primary:")
	positionPrim = tview.NewDropDown().SetLabel("Position:")
	relativeToPrim = tview.NewDropDown().SetLabel("Relative to:")
	saveButton = tview.NewButton("Save")
)

var (
	displays = GetDisplays()
	selectedDisplay = ""
	selectedDisplayIndex = 0
	selectedMode = "unchanged"
	position = " "
	relativeTo = " "
)

func setDisplays() {
	displayPrim.SetOptions(nil, nil)
	for _, display := range displays {
		displayPrim.AddOption(display.Name, func() {
			selectedDisplayIndex, selectedDisplay = displayPrim.GetCurrentOption()
			selectedDisplay := displays[selectedDisplayIndex]
			onPrim.SetChecked(!selectedDisplay.Off).SetChangedFunc(func(bool){update()})
			setModes()
			primaryPrim.SetChecked(selectedDisplay.Primary).SetChangedFunc(func(bool){update()})
			setPositions()
			setRelativeTo()
			update()
		})
	}
	displayPrim.SetCurrentOption(0)
}

func setModes() {
	modePrim.SetOptions(nil, nil)
	modes := append([]string{"unchanged", "auto"}, GetDisplayModes(selectedDisplay)...)
	for _, mode := range modes {
		modePrim.AddOption(mode, func() {
			_, selectedMode = modePrim.GetCurrentOption()
			update()
		})
	}
	modePrim.SetCurrentOption(0)
}

func setPositions() {
	var options []string
	if onlyDisplay(selectedDisplay) {
		options = []string{" "}
	} else {
		options = []string{" ", "same-as", "left-of", "right-of", "above", "below"}
	}

	positionPrim.SetOptions(
		options,
		func(pos string, _ int) {
			position = pos
			update()
		},
	)
	positionPrim.SetCurrentOption(0)
}

func setRelativeTo() {
	if onlyDisplay(selectedDisplay) {
		relativeToPrim.SetOptions([]string{" "}, func(string, int) {
			relativeTo = " "
			update()
		})
		relativeToPrim.SetCurrentOption(0)
		return
	}

	relativeToPrim.SetOptions(nil, nil)
	for _, display := range displays {
		if display.Name != selectedDisplay && !display.Off {
			relativeToPrim.AddOption(display.Name, func() {
				_, relativeTo = relativeToPrim.GetCurrentOption()
				update()
			})
		}
	}
	relativeToPrim.SetCurrentOption(0)
}

func onlyDisplay(name string) bool {
	for _, display := range displays {
		if !display.Off && display.Name != name {
			return false
		}
	}
	return true
}

func update() {
	command := getCommand()

	displayCopy := Copy(displays)
	ParseChanges(displayCopy, command)

	preview.SetText(DrawDisplays(displayCopy, 49, 25, 2).ToString())
}

func getCommand() string {
	displayArgument := "--output " + selectedDisplay

	var onArgument string
	if onPrim.IsChecked() {
		onArgument = ""
	} else {
		onArgument = "--off"
	}

	var modeArgument string
	switch (selectedMode) {
	case "unchanged":
		modeArgument = ""
	case "auto":
		modeArgument = "--auto"
	default:
		modeArgument = "--mode " + selectedMode
	}

	var primaryArgument string
	if primaryPrim.IsChecked() {
		primaryArgument = "--primary"
	} else {
		primaryArgument = "--noprimary"
	}

	var positionArgument string
	if position == " " {
		positionArgument = ""
	} else {
		positionArgument = "--" + position + " " + relativeTo
	}

	return fmt.Sprintf(
		`xrandr %s %s %s %s %s`,
		displayArgument, onArgument, modeArgument, primaryArgument, positionArgument,
	)
}

func save() {
	exec.Command("sh", "-c", getCommand()).Run()
	displays = GetDisplays()
	setDisplays()
}

func main() {
	preview = tview.NewTextView().SetWrap(false)

	form := tview.NewForm().
		AddFormItem(displayPrim).
		AddFormItem(onPrim).
		AddFormItem(modePrim).
		AddFormItem(primaryPrim).
		AddFormItem(positionPrim).
		AddFormItem(relativeToPrim).
		AddButton("Save", save)

	setDisplays()


	flex := tview.NewFlex().
		AddItem(preview, 50, 1, false).
		AddItem(form, 0, 1, true)

	if err := tview.NewApplication().SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
