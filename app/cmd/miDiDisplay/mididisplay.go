package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	lm "github.com/sirion/gomidi/lib/launchpadmini"
)

func main() {

	/**
	--buttons="a2:r3g2,c5:r2g1"
	--buttons="a2:AmberFull,c5:Yellow"
	*/

	device := flag.String("device", "auto", "Midi device for Launchpad")
	buttonStr := flag.String(
		"buttons",
		"",
		"List of comma separated buttons and their color values to set on the Launchpad.\n"+
			"Buttons and color values are separated by \":\".\n"+
			"Valid buttons are A-H, A1-H8 and L1-L8\n"+
			"Colors are described as their red and green values: r0g3 means full green, r3g0 means full red\n"+
			"Alternatively supported color names:\n"+
			"    ColorOff,\n"+
			"    ColorRedLow, ColorRedFull,\n"+
			"    ColorGreenLow, ColorGreenFull,\n"+
			"    ColorAmberLow, ColorAmberFull,\n"+
			"    ColorYellowFull\n"+
			"\n"+
			"Examples:\n"+
			"  a2:AmberFull - sets button A2 to the color amber with full brightness\n"+
			"  g3:r3g0      - sets button G3 to the color red with full brightness",
	)
	clear := flag.Bool("clear", false, "Start with all buttons turned off")
	help := flag.Bool("help", false, "Show this help")
	flag.Parse()

	lp := lm.New(*device)

	showHelp := *help || flag.NFlag() == 0
	if showHelp {
		fmt.Printf("Available command line options:\n")
		flag.PrintDefaults()
		return
	}

	if *clear {
		// TODO: If flashing is used reset cannot be used to clear the buttons without disabling flashing
		lp.Reset()
	}

	if *buttonStr != "" {
		setButtonsFromString(lp, *buttonStr)
	}

}

func setButtonsFromString(lp *lm.LaunchpadMini, buttonStr string) {

	buttons := strings.Split(buttonStr, ",")

	for _, button := range buttons {
		parts := strings.Split(button, ":")

		if len(parts) != 2 || len(parts[0]) > 2 || len(parts[1]) < 4 {
			fmt.Printf("Invalid button \"%s\", must have the format X:Y with X being the button name and Y being the color\n", button)
			continue
		}

		buttonName := strings.ToUpper(parts[0])
		var buttonID string
		if buttonName[0:1] == "L" {
			buttonID = "LiveButton" + buttonName[1:2]
		} else {
			buttonID = "Button" + buttonName
		}

		buttonValue, ok := lm.ButtonValues[buttonID]

		if !ok {
			fmt.Printf("Invalid button name \"%s\". Valid Names are A-H, A1-H8, L1-L8\n", buttonName)
			continue
		}

		color := parts[1]
		colorValue, ok := lm.ColorNames["Color"+color]
		if !ok {
			color := strings.ToUpper(color)

			r, err1 := strconv.ParseInt(color[1:2], 10, 8)
			g, err2 := strconv.ParseInt(color[3:4], 10, 8)

			if color[0:1] != "R" || color[2:3] != "G" || err1 != nil || err2 != nil || r > 3 || g > 3 || r < 0 || g < 0 {
				fmt.Printf("Invalid color value \"%s\" for button \"%s\". Please provide either a valid name or a value in the form of rXgY, with X and Y between 0 and 3\n", parts[1], parts[0])
				continue
			}

			colorValue = byte((0x10 * g) + 12 + r)
		}

		lp.Button(buttonValue, colorValue)

	}

}
