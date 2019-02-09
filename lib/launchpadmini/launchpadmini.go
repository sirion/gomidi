package launchpadmini

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sirion/gomidi/lib/midi"
)

// LaunchpadMini is the structure to connect to the Launchpad Mini midi device
type LaunchpadMini struct {
	fd        *os.File
	input     chan byte
	listening bool
}

// New creates a new instance of the LaunchpadMini and opens a connection to the given device
func New(device string) *LaunchpadMini {
	l := &LaunchpadMini{}

	if device == "auto" {
		device = findMidiDevice()
	}

	fd, err := os.OpenFile(device, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	l.fd = fd

	return l
}

func findMidiDevice() string {
	cmd := exec.Command("amidi", "-l")
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error listing USB devices: %s\n", err.Error())
	}

	list := strings.Split(string(stdout), "\n")

	if len(list) < 2 {
		log.Fatalf("Error finding USB device. %d devices found\n", len(list)-1)
	}

	var info string
	for _, line := range list {
		if strings.Contains(line, "Launchpad Mini") {
			info = line
			break
		}
	}

	if info == "" {
		log.Fatalf("Error finding USB device. Found %d midi devices. No Launchpad Mini\n", len(list)-1)
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
		log.Fatalf("Error parsing amidi list. hw-String did not contain three numbers: \"%s\"\n", info[3:])
	}
	return "/dev/snd/midiC" + parts[0] + "D" + parts[1]
}

// Button sets the given Button (from the Button* and LiveButton* constants) to the given color from the Color* constants
func (l *LaunchpadMini) Button(button, color byte) {
	if button < 204 {
		l.fd.Write(midi.NoteOn(0, button, color))
	} else {
		l.fd.Write(midi.Controller(0, button-100, color))
	}
	l.fd.Sync()
}

// Grid sets one of the grid buttons identified by its rown and column to the given color from the Color* constants
func (l *LaunchpadMini) Grid(row, column, color byte) {
	l.fd.Write(midi.NoteOn(0, (16*row)+column, color))
	l.fd.Sync()
}

// Live sets one of the live buttons identified by its number or the LiveButton* constant to the given color from the Color* constants
func (l *LaunchpadMini) Live(number, color byte) {
	l.fd.Write(midi.Controller(0, 104+number, color))
	l.fd.Sync()
}

// Reset sets all buttons to off and clears all other settings made in the session
func (l *LaunchpadMini) Reset() {
	l.fd.Write([]byte{176, 0, 0})
	l.fd.Sync()
}

// Listen returns a channel containing the pressed keys on the launchpad.
// The key is sent to the channel on button down only and contains the key as described by the Button* and LiveButton* constants
func (l *LaunchpadMini) Listen() chan byte {
	l.input = make(chan byte, 1)
	l.listening = true

	go func() {
		buffer := make([]byte, 3, 3)
		for l.listening {
			read, err := l.fd.Read(buffer)
			if err != nil {
				fmt.Printf("Reading error from Launchpad: %s\n", err.Error())
				return
			}

			if read != 3 || buffer[2] != 127 {
				if buffer[2] != 0 {
					fmt.Printf("Unknown Launchpad input: %#v\n", buffer[:read])
				}
				continue
			}

			if buffer[0] == 144 {
				// Grid Button
				l.input <- buffer[1]
			} else if buffer[0] == 176 {
				// Live Button
				l.input <- buffer[1] + 100
			} else {
				// Ignored
				fmt.Printf("Unknown: %d (%d)\n", buffer[1], buffer[2])
			}
		}

	}()

	return l.input
}

// Text outputs a string to the launchpad in the given color from the COlor* constants
func (l *LaunchpadMini) Text(text string, color byte) {
	l.fd.Write(textCmdPrefix)
	l.fd.Write([]byte{color})
	l.fd.Write([]byte(text))
	l.fd.Write(textCmdSuffix)
	l.fd.Sync()
}

// AllOn sets all LEDs to amber with the given intensity between 125 and 127
func (l *LaunchpadMini) AllOn(intensity byte) {
	if intensity < 125 {
		intensity = 125
	} else if intensity > 127 {
		intensity = 127
	}

	l.fd.Write([]byte{176, 0, intensity})
	l.fd.Sync()
}

// Flashing turns on flashing buttons (at a default speed). Cannot be used with double buffering at the same time
func (l *LaunchpadMini) Flashing(on bool) {
	if on {
		l.fd.Write([]byte{0xb0, 0, 0x28})
	} else {
		l.fd.Write([]byte{0xb0, 0, 0x30})
	}
	l.fd.Sync()
}

// RapidUpdate sets the LED status of all launchpad buttons at once.
// The given buttonmap contains all values. Buttons not set in the map will be set to off.
func (l *LaunchpadMini) RapidUpdate(buttonmap map[byte]byte) {
	var x, y, i byte
	buttons := make([]byte, 81, 81)

	// Send a NoteOn on Channel 3, then 80 colors. Button order: Grid, A-H, Live
	buttons[i] = 0x92
	i++

	// In order to improve readability I will not spell out all buttons and instead loop over the grid.
	for y = 0; y < 8; y++ {
		for x = 0; x < 8; x++ {
			buttons[i] = buttonmap[x+16*y]
			i++
		}
	}

	buttons[i] = buttonmap[ButtonA]
	i++
	buttons[i] = buttonmap[ButtonB]
	i++
	buttons[i] = buttonmap[ButtonC]
	i++
	buttons[i] = buttonmap[ButtonD]
	i++
	buttons[i] = buttonmap[ButtonE]
	i++
	buttons[i] = buttonmap[ButtonF]
	i++
	buttons[i] = buttonmap[ButtonG]
	i++
	buttons[i] = buttonmap[ButtonH]
	i++

	buttons[i] = buttonmap[LiveButton1]
	i++
	buttons[i] = buttonmap[LiveButton2]
	i++
	buttons[i] = buttonmap[LiveButton3]
	i++
	buttons[i] = buttonmap[LiveButton4]
	i++
	buttons[i] = buttonmap[LiveButton5]
	i++
	buttons[i] = buttonmap[LiveButton6]
	i++
	buttons[i] = buttonmap[LiveButton7]
	i++
	buttons[i] = buttonmap[LiveButton8]

	l.fd.Write(buttons)
	l.fd.Sync()
}

// BufferMode sets the working mode of the launchpad. The most useful values are provided as BufferMode*-constants
//
// The mode byte works like this:
//
// Byte b01234567
//  0: always 0 -> data byte
//  1: always 0
//  2: always 1
//  3: Copy     -> copy from display to update buffer
//  4: Flash    -> flip buffers to make led flash
//  5: Update   -> set buffer as the new updating buffer
//  6: always 0
//  7: Display  -> set buffer as the new displaying buffer
//
// Examples:
//  b00100100 --> 0x24: Set buffer 1 to update
//  b00100001 --> 0x21: Set buffer 0 to update
//  b00101000 --> 0x28: Set mode to flashing
//  b00110000 --> 0x30: Default mode: Update both buffers (double buffering off)
//
func (l *LaunchpadMini) BufferMode(mode byte) {
	// Check masks:
	//  b00100000 --> 0x20: mode & 0x20 == 0x20 --> Check if bit 2 is set to 1
	//  b11000010 --> 0xc2: mode & 0xc2 == 0x00 --> Check if bits 0, 1 and 6 are set to 0
	//  b00111101 --> 0x3d: mode & 0x3d --> only allow 1s for bits 2, 3, 4, 5 and 7
	//
	// Make sure bit 2 is set to 1 and bits 0, 1 and 6 are set to 0
	mode = (mode | 0x20) & 0x3d

	// Todo: What happens when bit 5 and 7 are set to 1?

	l.fd.Write([]byte{0xb0, 0, mode})
	l.fd.Sync()
}

// Close closes the connection to the midi device and optionally ends the listening process
func (l *LaunchpadMini) Close() {
	l.listening = false
	l.fd.Close()
}
