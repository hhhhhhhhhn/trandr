package main

import (
	"os/exec"
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

type Display struct {
	Name string
	Off bool
	Width int
	Height int
	X int
	Y int
	Primary bool
}

var (
	modeRegex = regexp.MustCompile(`\d*x\d*\+\d*\+\d*`)
	primaryRegex = regexp.MustCompile(` primary `)

	nonDigitRegex = regexp.MustCompile(`\D+`)

	displayNameRegex = regexp.MustCompile(`"[^"]*"`)
	resolutionRegex = regexp.MustCompile(`\d*x\d*`)
	positionRegex = regexp.MustCompile(`\+\d*\+\d*`)
)

func GetDisplays() (diplays []Display) {
	out, _ := exec.Command("sh", "-c", "xrandr | grep ' connected'").Output()

	for _, line := range strings.Split(string(out[:len(out) - 1]), "\n") {
		diplays = append(diplays, ParseDisplay(line))
		fmt.Println(line, ParseDisplay(line))
	}
	return diplays
}

func ParseDisplay(line string) Display {
	display := Display{}
	display.Name = strings.Split(line, " ")[0]
	modeString := modeRegex.FindString(line)
	if modeString == "" {
		display.Off = true
	} else {
		display.ParseMode(modeString)
	}
	display.Primary = primaryRegex.MatchString(line)
	return display
}

// in wxh+x+y
func (d* Display) ParseMode(mode string) {
	fields := nonDigitRegex.Split(mode, -1)

	d.Width, _  = strconv.Atoi(fields[0])
	d.Height, _ = strconv.Atoi(fields[1])
	d.X, _      = strconv.Atoi(fields[2])
	d.Y, _      = strconv.Atoi(fields[3])
}

func DrawDisplays (diplays []Display, width, height, aspectRatio int) Grid {
	g := NewGrid(width, height)
	maxX, maxY := getMaxXAndY(diplays)
	scalingFactor := getScalingFactor(maxX * aspectRatio, maxY, width, height)

	for _, display := range diplays {
		if !display.Off {
			g.DrawTextBox(
				display.Name,
				display.X * aspectRatio / scalingFactor,
				display.Y / scalingFactor,
				(display.X + display.Width) * aspectRatio / scalingFactor,
				(display.Y + display.Height) / scalingFactor,
			)
		}
	}

	return g
}

func getMaxXAndY(displays []Display) (int, int) {
	maxX := 0
	maxY := 0

	for _, display := range displays {
		if display.X + display.Width > maxX {
			maxX = display.X + display.Width
		}
		if display.Y + display.Height > maxY {
			maxY = display.Y + display.Height
		}
	}
	return maxX, maxY
}

func getScalingFactor(maxX, maxY, width, height int) int {
	if maxX * 1000 / width > maxY * 1000 / height {
		return maxX / (width - 2)
	} else {
		return maxY / (height - 2)
	}
}

func Copy(displays []Display) []Display {
	displaysCopy := []Display{}

	for _, display := range displays {
		displayCopy := Display{
			Name: display.Name,
			Off: display.Off,
			Width: display.Width,
			Height: display.Height,
			X: display.X,
			Y: display.Y,
			Primary: display.Primary,
		}
		displaysCopy = append(displaysCopy, displayCopy)
	}

	return displaysCopy
}

func ParseChanges(displays []Display, command string) {
	out, _ := exec.Command("sh", "-c", command + " --dryrun | grep '\"'").Output()

	if len(out) == 0 {
		return
	}

	for _, line := range strings.Split(string(out[:len(out) - 1]), "\n") {
		ParseChange(displays, line)
	}
}

func ParseChange(displays []Display, line string) {
	name := displayNameRegex.FindString(line)
	name = name[1:len(name)-1] // gets rid of quotes
	modeString := resolutionRegex.FindString(line) + positionRegex.FindString(line)

	for i := range displays {
		if displays[i].Name == name {
			displays[i].Off = false
			displays[i].ParseMode(modeString)
		}
	}
}

func GetDisplayModes(displayName string) []string {
	// The first sed expression deletes up to the wanted display
	// The second removes everything after the modes (uses the fact they're indented)
	// The third one selects the modes out of the rest of the line
	command := fmt.Sprintf(
		`xrandr | sed -E -e '1,/%s/ d' -e '/^\S/,$ d' -e 's/\s*(\S*).*/\1/'`,
		displayName,
	)
	out, _ := exec.Command("sh", "-c", command).Output()

	return strings.Split(string(out[:len(out)-1]), "\n")
}


//func main() {
//	// TODO: TUI d := GetDisplays()
//
//	c := Copy(d)
//
//	ParseChanges(c, "xrandr --output eDP1 --mode 1280x720")
//
//	fmt.Println(d, c)
//	fmt.Println(GetDisplayModes(d[0].Name))
//
//	fmt.Println(DrawDisplays(d, 50, 15, 2).ToString())
//	fmt.Println(DrawDisplays(c, 50, 15, 2).ToString())
//}
