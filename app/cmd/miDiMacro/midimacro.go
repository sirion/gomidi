package main

/**
 * ALMOST:
 *  - Automatically find midi device if not given/configured or set to "auto"
 *    TODO: make sure this actually works every time
 *
 * TODOs:
 *
 *	- Support more than just Keyboard shortcuts (maybe mouse macros?)
 *  - Port to Windows :-/
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
	"os/user"
	"path/filepath"

	lm "github.com/sirion/gomidi/lib/launchpadmini"

	"github.com/go-vgo/robotgo"
)

// KeyCombination describes the macro key combination and consists of one key and optional modifiers
type KeyCombination struct {
	Key       string   `json:"key"`
	Modifiers []string `json:"modifiers,omitempty"`
}

type Configuration struct {
	Macros map[byte]KeyCombination `json:"-"`

	Device    string                    `json:"device"`
	KeyMacros map[string]KeyCombination `json:"keyMacros"`
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

	c.KeyMacros = make(map[string]KeyCombination, len(c.Macros))
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

	if c.Device == "" {
		c.Device = "auto"
	}

	c.Macros = make(map[byte]KeyCombination, len(c.KeyMacros))
	for key, combo := range c.KeyMacros {
		c.Macros[lm.ButtonValues[key]] = combo
	}
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
func pressKey(kc KeyCombination) {
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

//var keyMap = map[byte]KeyCombination{
// 	lm.LiveButton1: KeyCombination{Key: "1", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton2: KeyCombination{Key: "2", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton3: KeyCombination{Key: "3", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton4: KeyCombination{Key: "4", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton5: KeyCombination{Key: "5", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton6: KeyCombination{Key: "6", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton7: KeyCombination{Key: "7", Modifiers: []string{"ctrl", "alt"}},
// 	lm.LiveButton8: KeyCombination{Key: "8", Modifiers: []string{"ctrl", "alt"}},

// 	lm.ButtonA1: KeyCombination{Key: "1", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA2: KeyCombination{Key: "2", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA3: KeyCombination{Key: "3", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA4: KeyCombination{Key: "4", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA5: KeyCombination{Key: "5", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA6: KeyCombination{Key: "6", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA7: KeyCombination{Key: "7", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA8: KeyCombination{Key: "8", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonA:  KeyCombination{Key: "9", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB1: KeyCombination{Key: "0", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB2: KeyCombination{Key: "a", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB3: KeyCombination{Key: "b", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB4: KeyCombination{Key: "c", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB5: KeyCombination{Key: "d", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB6: KeyCombination{Key: "e", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB7: KeyCombination{Key: "f", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB8: KeyCombination{Key: "g", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonB:  KeyCombination{Key: "h", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC1: KeyCombination{Key: "i", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC2: KeyCombination{Key: "j", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC3: KeyCombination{Key: "k", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC4: KeyCombination{Key: "l", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC5: KeyCombination{Key: "m", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC6: KeyCombination{Key: "n", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC7: KeyCombination{Key: "o", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC8: KeyCombination{Key: "p", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonC:  KeyCombination{Key: "q", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD1: KeyCombination{Key: "r", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD2: KeyCombination{Key: "s", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD3: KeyCombination{Key: "t", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD4: KeyCombination{Key: "u", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD5: KeyCombination{Key: "v", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD6: KeyCombination{Key: "w", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD7: KeyCombination{Key: "x", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD8: KeyCombination{Key: "y", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonD:  KeyCombination{Key: "z", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonE1: KeyCombination{Key: "num0", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE2: KeyCombination{Key: "num1", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE3: KeyCombination{Key: "num2", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE4: KeyCombination{Key: "num3", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE5: KeyCombination{Key: "num4", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE6: KeyCombination{Key: "num5", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE7: KeyCombination{Key: "num6", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE8: KeyCombination{Key: "num7", Modifiers: []string{"ctrl"}},
// 	lm.ButtonE:  KeyCombination{Key: "num8", Modifiers: []string{"ctrl"}},
// 	lm.ButtonF1: KeyCombination{Key: "num9", Modifiers: []string{"ctrl"}},
// 	lm.ButtonF2: KeyCombination{Key: "num0", Modifiers: []string{"alt"}},
// 	lm.ButtonF3: KeyCombination{Key: "num1", Modifiers: []string{"alt"}},
// 	lm.ButtonF4: KeyCombination{Key: "num2", Modifiers: []string{"alt"}},
// 	lm.ButtonF5: KeyCombination{Key: "num3", Modifiers: []string{"alt"}},
// 	lm.ButtonF6: KeyCombination{Key: "num4", Modifiers: []string{"alt"}},
// 	lm.ButtonF7: KeyCombination{Key: "num5", Modifiers: []string{"alt"}},
// 	lm.ButtonF8: KeyCombination{Key: "num6", Modifiers: []string{"alt"}},
// 	lm.ButtonF:  KeyCombination{Key: "num7", Modifiers: []string{"alt"}},
// 	lm.ButtonG1: KeyCombination{Key: "num8", Modifiers: []string{"alt"}},
// 	lm.ButtonG2: KeyCombination{Key: "num9", Modifiers: []string{"alt"}},
// 	lm.ButtonG3: KeyCombination{Key: "num0", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG4: KeyCombination{Key: "num1", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG5: KeyCombination{Key: "num2", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG6: KeyCombination{Key: "num3", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG7: KeyCombination{Key: "num4", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG8: KeyCombination{Key: "num5", Modifiers: []string{"ctrl"}},
// 	lm.ButtonG:  KeyCombination{Key: "num6", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH1: KeyCombination{Key: "num7", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH2: KeyCombination{Key: "num8", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH3: KeyCombination{Key: "num9", Modifiers: []string{"ctrl"}},
// 	lm.ButtonH4: KeyCombination{Key: "num0", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH5: KeyCombination{Key: "num1", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH6: KeyCombination{Key: "num2", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH7: KeyCombination{Key: "num3", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH8: KeyCombination{Key: "num4", Modifiers: []string{"ctrl", "cmd"}},
// 	lm.ButtonH:  KeyCombination{Key: "num5", Modifiers: []string{"ctrl", "cmd"}},
// }
