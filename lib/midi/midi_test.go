package midi

import (
	"fmt"
	"strings"
	"testing"
)

func byteSliceToInterfaceSlice(byteSlc []byte) []interface{} {
	iSlc := make([]interface{}, len(byteSlc), len(byteSlc))
	for i, b := range byteSlc {
		iSlc[i] = interface{}(b)
	}
	return iSlc
}

func compareBytes(t *testing.T, got, want []byte) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Wrong Size: Got %d, expected %d", len(got), len(want))
	}

	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {

			gotString := fmt.Sprintf(strings.Repeat("%x ", len(got)), byteSliceToInterfaceSlice(got)...)
			wantString := fmt.Sprintf(strings.Repeat("%x ", len(want)), byteSliceToInterfaceSlice(want)...)

			t.Errorf("Wrong value at position %d: Got { %s}, expected { %s}", i, gotString, wantString)
		}
	}
}

/*
 *  Note on:     1001cccc 0ppppppp 0vvvvvvv (c = channel, p = pitch, v = velocity)
 */
func TestNoteOn(t *testing.T) {
	compareBytes(t, NoteOn(0, 127, 127), []byte{0x90, 127, 127})
	compareBytes(t, NoteOn(0, 0, 0), []byte{0x90, 0, 0})
	compareBytes(t, NoteOn(0, 77, 93), []byte{0x90, 77, 93})

	compareBytes(t, NoteOn(1, 0, 0), []byte{0x91, 0, 0})
	compareBytes(t, NoteOn(10, 0, 0), []byte{0x9a, 0, 0})
	compareBytes(t, NoteOn(15, 0, 0), []byte{0x9f, 0, 0})
	compareBytes(t, NoteOn(16, 0, 0), []byte{0x9f, 0, 0})
	compareBytes(t, NoteOn(175, 0, 0), []byte{0x9f, 0, 0})

	// Max value for pitch and velocity is 127
	compareBytes(t, NoteOn(0, 255, 255), []byte{0x90, 127, 127})
	compareBytes(t, NoteOn(0, 231, 243), []byte{0x90, 127, 127})
}

/*
 *  Note off:    1000cccc 0ppppppp 0vvvvvvv (c = channel, p = pitch, v = velocity) -> v is usually 0
 */
func TestNoteOff(t *testing.T) {
	compareBytes(t, NoteOff(0, 127, 127), []byte{0x80, 127, 127})
	compareBytes(t, NoteOff(0, 0, 0), []byte{0x80, 0, 0})
	compareBytes(t, NoteOff(0, 77, 93), []byte{0x80, 77, 93})

	compareBytes(t, NoteOff(1, 0, 0), []byte{0x81, 0, 0})
	compareBytes(t, NoteOff(10, 0, 0), []byte{0x8a, 0, 0})
	compareBytes(t, NoteOff(15, 0, 0), []byte{0x8f, 0, 0})
	compareBytes(t, NoteOff(16, 0, 0), []byte{0x8f, 0, 0})
	compareBytes(t, NoteOff(175, 0, 0), []byte{0x8f, 0, 0})

	// Max value for pitch and velocity is 127
	compareBytes(t, NoteOff(0, 255, 255), []byte{0x80, 127, 127})
	compareBytes(t, NoteOff(0, 231, 243), []byte{0x80, 127, 127})
}

func TestDrumOn(t *testing.T) {
	// Drums are just nores on channel 10 (9)

	compareBytes(t, DrumOn(127, 127), []byte{0x99, 127, 127})
	compareBytes(t, DrumOn(0, 0), []byte{0x99, 0, 0})
	compareBytes(t, DrumOn(77, 93), []byte{0x99, 77, 93})

	// Max value for pitch and velocity is 127
	compareBytes(t, DrumOn(255, 255), []byte{0x99, 127, 127})
	compareBytes(t, DrumOn(231, 243), []byte{0x99, 127, 127})
}

func TestDrumOff(t *testing.T) {
	// Drums are just nores on channel 10 (9)

	compareBytes(t, DrumOff(127, 127), []byte{0x89, 127, 127})
	compareBytes(t, DrumOff(0, 0), []byte{0x89, 0, 0})
	compareBytes(t, DrumOff(77, 93), []byte{0x89, 77, 93})

	// Max value for pitch and velocity is 127
	compareBytes(t, DrumOff(255, 255), []byte{0x89, 127, 127})
	compareBytes(t, DrumOff(231, 243), []byte{0x89, 127, 127})
}

/*
 *  Controller:  1011cccc 0nnnnnnn 0vvvvvvv (c = channel, n = controller number, v = value)
 */
func TestController(t *testing.T) {
	compareBytes(t, Controller(0, 127, 127), []byte{0xb0, 127, 127})
	compareBytes(t, Controller(0, 0, 0), []byte{0xb0, 0, 0})
	compareBytes(t, Controller(0, 77, 93), []byte{0xb0, 77, 93})

	compareBytes(t, Controller(1, 0, 0), []byte{0xb1, 0, 0})
	compareBytes(t, Controller(10, 0, 0), []byte{0xba, 0, 0})
	compareBytes(t, Controller(15, 0, 0), []byte{0xbf, 0, 0})
	compareBytes(t, Controller(16, 0, 0), []byte{0xbf, 0, 0})
	compareBytes(t, Controller(175, 0, 0), []byte{0xbf, 0, 0})

	// Max value for pitch and velocity is 127
	compareBytes(t, Controller(0, 255, 255), []byte{0xb0, 127, 127})
	compareBytes(t, Controller(0, 231, 243), []byte{0xb0, 127, 127})
}

/*
 *  Prog Change: 1100cccc 0xxxxxxx          (c = channel, x = instrument number)
 */
func TestProgramChange(t *testing.T) {
	compareBytes(t, ProgramChange(0, 127), []byte{0xc0, 127})
	compareBytes(t, ProgramChange(0, 0), []byte{0xc0, 0})
	compareBytes(t, ProgramChange(0, 93), []byte{0xc0, 93})

	compareBytes(t, ProgramChange(1, 0), []byte{0xc1, 0})
	compareBytes(t, ProgramChange(10, 0), []byte{0xca, 0})
	compareBytes(t, ProgramChange(15, 0), []byte{0xcf, 0})
	compareBytes(t, ProgramChange(16, 0), []byte{0xcf, 0})
	compareBytes(t, ProgramChange(175, 0), []byte{0xcf, 0})

	// Max value for pitch and velocity is 127
	compareBytes(t, ProgramChange(0, 255), []byte{0xc0, 127})
	compareBytes(t, ProgramChange(0, 243), []byte{0xc0, 127})
}

/*
 *  Bend Pitch:  1110cccc 0vvvvvvv 0vvvvvvv (c = channel, v = value) -> v has 14 bit 0 to 16383 (0x3fff)
 */
func TestBendPitch(t *testing.T) {
	compareBytes(t, BendPitch(0, 127), []byte{0xe0, 0, 127})
	compareBytes(t, BendPitch(0, 0), []byte{0xe0, 0, 0})
	compareBytes(t, BendPitch(0, 93), []byte{0xe0, 0, 93})
	compareBytes(t, BendPitch(0, 128), []byte{0xe0, 1, 0})
	compareBytes(t, BendPitch(0, 129), []byte{0xe0, 1, 1})

	compareBytes(t, BendPitch(1, 0), []byte{0xe1, 0, 0})
	compareBytes(t, BendPitch(10, 0), []byte{0xea, 0, 0})
	compareBytes(t, BendPitch(15, 0), []byte{0xef, 0, 0})
	compareBytes(t, BendPitch(16, 0), []byte{0xef, 0, 0})
	compareBytes(t, BendPitch(175, 0), []byte{0xef, 0, 0})

	// Max value for value is 16383
	compareBytes(t, BendPitch(0, 16384), []byte{0xe0, 127, 127})
	compareBytes(t, BendPitch(0, 43214), []byte{0xe0, 127, 127})
}

/*
 *  Reset:       11111111
 */
func TestReset(t *testing.T) {
	if Reset()[0] != 0xff {
		t.Fail()
	}

	if Reset()[0] != b11111111 {
		t.Fail()
	}

	if len(Reset()) != 1 {
		t.Fail()
	}
}
