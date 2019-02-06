package launchpadmini

import (
	"fmt"
	"os"

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

	fd, err := os.OpenFile(device, os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	l.fd = fd

	l.Reset()

	return l
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
				fmt.Printf("Reading error from Launchpad: %s", err.Error())
				return
			}

			if read != 3 || buffer[2] != 127 {
				if buffer[2] != 0 {
					fmt.Printf("Unknown Launchpad input: %#v", buffer[:read])
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

// Close closes the connection to the midi device and optionally ends the listening process
func (l *LaunchpadMini) Close() {
	l.listening = false
	l.fd.Close()
}
