package main

/**
 * ALMOST:
 *  - Automatically find midi device if not given/configured or set to "auto"
 *    TODO: make sure this actually works every time
 *
 * TODOs:
 *  - Split into two tools: 1. Read and convert buttons. 2. Set button states from command line
 *
 *	- Support more than just Keyboard shortcuts (maybe mouse macros?)
 *  - Create documentation
 *  - Port to Windows :-/
 *
 *  - Port to / test on Mac :-/
 *
 */

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	lm "github.com/sirion/gomidi/lib/launchpadmini"

	"github.com/go-vgo/robotgo"
	// "github.com/kylelemons/gousb/usb"
)

type Configuration struct {
	Macros map[byte]lm.KeyCombination `json:"-"`

	Device    string                       `json:"device"`
	KeyMacros map[string]lm.KeyCombination `json:"keyMacros"`
}

func getUserDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("Could not find user name: %s", err.Error())
	}

	return user.HomeDir
}

func (c *Configuration) Save() {
	directory := filepath.Join(getUserDir(), ".config", "midimacro")
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		log.Fatalf("Error saving configuration file: %s", err.Error())
	}

	c.KeyMacros = make(map[string]lm.KeyCombination, len(c.Macros))
	for key, combo := range c.Macros {
		c.KeyMacros[lm.ButtonNames[key]] = combo
	}

	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting configuration file: %s", err.Error())
	}

	ioutil.WriteFile(filepath.Join(directory, "config.json"), out, os.ModePerm)
}

func (c *Configuration) Load() {
	device := flag.String("device", "", "Override configured midi device path")
	configurationPath := flag.String("config", "", "Override configuration file path")
	flag.Parse()

	if *configurationPath == "" {
		directory := filepath.Join(getUserDir(), ".config", "midimacro")
		*configurationPath = filepath.Join(directory, "config.json")
	}

	data, err := ioutil.ReadFile(*configurationPath)
	if err != nil {
		log.Fatalf("Error loading configuration file: %s", err.Error())
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %s", err.Error())
	}

	if *device != "" {
		c.Device = *device
	}

	if c.Device == "auto" {
		c.FindMidiDevice()
	}

	c.Macros = make(map[byte]lm.KeyCombination, len(c.KeyMacros))
	for key, combo := range c.KeyMacros {
		c.Macros[lm.ButtonValues[key]] = combo
	}
}

func (c *Configuration) FindMidiDevice() {
	cmd := exec.Command("amidi", "-l")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error listing USB devices: %s", err.Error())
	}

	list := strings.Split(string(stdout), "\n")

	if len(list) < 2 {
		log.Fatalf("Error finding USB device. %d devices found", len(list)-1)
	}

	var info string
	for _, line := range list {
		if strings.Contains(line, "Launchpad Mini") {
			info = line
			break
		}
	}

	if info == "" {
		log.Fatalf("Error finding USB device. Found %d midi devices. No Launchpad Mini", len(list)-1)
	}

	parts := strings.Split(info, " ")
	i := 0
	for _, str := range parts {
		if str != "" {
			i++
		}
		if i == 2 {
			info = str
			break
		}
	}

	parts = strings.Split(info[3:], ",")

	if len(parts) != 3 {
		log.Fatalf("Error parsing amidi list. hw-String did not contain three numbers: \"%s\"", info[3:])
	}
	c.Device = "/dev/snd/midiC" + parts[0] + "D" + parts[1]
}

func main() {
	config := Configuration{}
	config.Load()

	lp := lm.New(config.Device)

	input := lp.Listen()

	for {
		press := <-input

		keys, ok := config.Macros[press]
		if ok {
			pressKey(keys)
		}
	}

}

// pressKey only exists because robotgo.KeyTap(string, []string...) does not work
func pressKey(kc lm.KeyCombination) {
	switch len(kc.Modifiers) {
	case 0:
		robotgo.KeyTap(kc.Key)
		break
	case 1:
		robotgo.KeyTap(kc.Key, kc.Modifiers[0])
		break
	case 2:
		robotgo.KeyTap(kc.Key, kc.Modifiers[0], kc.Modifiers[1])
		break
	case 3:
		robotgo.KeyTap(kc.Key, kc.Modifiers[0], kc.Modifiers[1], kc.Modifiers[2])
		break
	case 4:
		robotgo.KeyTap(kc.Key, kc.Modifiers[0], kc.Modifiers[1], kc.Modifiers[2], kc.Modifiers[3])
		break
	case 5:
		robotgo.KeyTap(kc.Key, kc.Modifiers[0], kc.Modifiers[1], kc.Modifiers[2], kc.Modifiers[3], kc.Modifiers[4])
		break
	default:
		fmt.Print("More than 5 modifiers are not supported")
	}
}

// const devicePath = "/dev/snd/midiC4D0"

//var keyMap = map[byte]lm.KeyCombination{
// 	lm.LiveButton1: lm.KeyCombination{Key: "1", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton2: lm.KeyCombination{Key: "2", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton3: lm.KeyCombination{Key: "3", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton4: lm.KeyCombination{Key: "4", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton5: lm.KeyCombination{Key: "5", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton6: lm.KeyCombination{Key: "6", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton7: lm.KeyCombination{Key: "7", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton8: lm.KeyCombination{Key: "8", Modifiers: []string{"ctrl", "alt"}},

// 	lm.ButtonA1: lm.KeyCombination{Key: "1", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA2: lm.KeyCombination{Key: "2", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA3: lm.KeyCombination{Key: "3", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA4: lm.KeyCombination{Key: "4", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA5: lm.KeyCombination{Key: "5", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA6: lm.KeyCombination{Key: "6", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA7: lm.KeyCombination{Key: "7", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA8: lm.KeyCombination{Key: "8", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA:  lm.KeyCombination{Key: "9", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB1: lm.KeyCombination{Key: "0", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB2: lm.KeyCombination{Key: "a", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB3: lm.KeyCombination{Key: "b", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB4: lm.KeyCombination{Key: "c", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB5: lm.KeyCombination{Key: "d", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB6: lm.KeyCombination{Key: "e", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB7: lm.KeyCombination{Key: "f", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB8: lm.KeyCombination{Key: "g", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB:  lm.KeyCombination{Key: "h", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC1: lm.KeyCombination{Key: "i", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC2: lm.KeyCombination{Key: "j", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC3: lm.KeyCombination{Key: "k", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC4: lm.KeyCombination{Key: "l", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC5: lm.KeyCombination{Key: "m", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC6: lm.KeyCombination{Key: "n", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC7: lm.KeyCombination{Key: "o", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC8: lm.KeyCombination{Key: "p", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC:  lm.KeyCombination{Key: "q", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD1: lm.KeyCombination{Key: "r", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD2: lm.KeyCombination{Key: "s", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD3: lm.KeyCombination{Key: "t", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD4: lm.KeyCombination{Key: "u", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD5: lm.KeyCombination{Key: "v", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD6: lm.KeyCombination{Key: "w", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD7: lm.KeyCombination{Key: "x", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD8: lm.KeyCombination{Key: "y", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD:  lm.KeyCombination{Key: "z", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonE1: lm.KeyCombination{Key: "num0", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE2: lm.KeyCombination{Key: "num1", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE3: lm.KeyCombination{Key: "num2", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE4: lm.KeyCombination{Key: "num3", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE5: lm.KeyCombination{Key: "num4", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE6: lm.KeyCombination{Key: "num5", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE7: lm.KeyCombination{Key: "num6", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE8: lm.KeyCombination{Key: "num7", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE:  lm.KeyCombination{Key: "num8", Modifiers: []string{"ctrl"}},
// 	lm.ButtonF1: lm.KeyCombination{Key: "num9", Modifiers: []string{"ctrl"}},
// 	lm.ButtonF2: lm.KeyCombination{Key: "num0", Modifiers: []string{"alt"}},
// 	lm.ButtonF3: lm.KeyCombination{Key: "num1", Modifiers: []string{"alt"}},
// 	lm.ButtonF4: lm.KeyCombination{Key: "num2", Modifiers: []string{"alt"}},
// 	lm.ButtonF5: lm.KeyCombination{Key: "num3", Modifiers: []string{"alt"}},
// 	lm.ButtonF6: lm.KeyCombination{Key: "num4", Modifiers: []string{"alt"}},
// 	lm.ButtonF7: lm.KeyCombination{Key: "num5", Modifiers: []string{"alt"}},
// 	lm.ButtonF8: lm.KeyCombination{Key: "num6", Modifiers: []string{"alt"}},
// 	lm.ButtonF:  lm.KeyCombination{Key: "num7", Modifiers: []string{"alt"}},
// 	lm.ButtonG1: lm.KeyCombination{Key: "num8", Modifiers: []string{"alt"}},
// 	lm.ButtonG2: lm.KeyCombination{Key: "num9", Modifiers: []string{"alt"}},
// 	lm.ButtonG3: lm.KeyCombination{Key: "num0", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG4: lm.KeyCombination{Key: "num1", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG5: lm.KeyCombination{Key: "num2", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG6: lm.KeyCombination{Key: "num3", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG7: lm.KeyCombination{Key: "num4", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG8: lm.KeyCombination{Key: "num5", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG:  lm.KeyCombination{Key: "num6", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH1: lm.KeyCombination{Key: "num7", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH2: lm.KeyCombination{Key: "num8", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH3: lm.KeyCombination{Key: "num9", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH4: lm.KeyCombination{Key: "num0", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH5: lm.KeyCombination{Key: "num1", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH6: lm.KeyCombination{Key: "num2", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH7: lm.KeyCombination{Key: "num3", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH8: lm.KeyCombination{Key: "num4", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH:  lm.KeyCombination{Key: "num5", Modifiers: []string{"ctrl", "cmd"}},
// }
