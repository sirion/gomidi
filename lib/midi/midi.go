// Package midi only contains helper functions that create midi-messages (byte-slices).
// It is essentially the code version of what I learned from reading http://www.music-software-development.com/midi-tutorial.html.
//
// Rule for all midi messages: The first bit in every byte determines whether it is a command byte (1) or a data byte (0) which follows a command byte
//
// Midi-Messages (binary notation):
//  Note off:    1000cccc 0ppppppp 0vvvvvvv (c = channel, p = pitch, v = velocity) -> v is usually 0
//  Note on:     1001cccc 0ppppppp 0vvvvvvv (c = channel, p = pitch, v = velocity)
//  Controller:  1011cccc 0nnnnnnn 0vvvvvvv (c = channel, n = controller number, v = value)
//  Prog Change: 1100cccc 0xxxxxxx          (c = channel, x = instrument number)
//  Bend Pitch:  1110cccc 0vvvvvvv 0vvvvvvv (c = channel, v = value) -> v has 14 bit 0 to 16383 (0x3fff)
//  Reset:       11111111
//
//  - LSB = Least Significant Byte, MSB = Most Significant Byte
//
package midi

//	Since go does not support binary notation, for enhanced readability,
//	here are the values needed for the midi masks:
//
//	Note off:    10000000: 0x80    10001111: 0x8f
//	Note on:     10010000: 0x90    10011111: 0x9f
//	             10100000: 0xa0    10101111: 0xaf
//	Controller:  10110000: 0xb0    10111111: 0xbf
//	Prog Change: 11000000: 0xc0    11001111: 0xcf
//	             11010000: 0xd0    11011111: 0xdf
//	Bend Pitch:  11100000: 0xe0    11101111: 0xef
//	             11110000: 0xf0    11111111: 0xff
//
//	General:
//	             11110000: 0xf0    00001111: 0x0f
//	             01111111: 0x7f    10000000: 0x80
//	             11111111: 0xff    00000000: 0x00
//
//	14Bit (for Pitch Bend):
//	             0011111111111111:  0x3fff
//
// The following constants are used to make the midi bit masks more readable
const (
	b10000000         = 0x80   // Note off
	b10010000         = 0x90   // Note on
	b10110000         = 0xb0   // Controller
	b11000000         = 0xc0   // Program Change
	b11100000         = 0xe0   // Bend Pitch
	b11111111         = 0xff   // Reset
	b01111111         = 0x7f   // Data mask / max value
	b00001111         = 0x0f   // Channel number mask
	b0011111111111111 = 0x3fff // 14bit for pitch bend value
)

// NoteOn sends the signal to start playing a note on the given
// channel using the given pitch and velocity
func NoteOn(channel, pitch, velocity byte) []byte {
	if channel > b00001111 {
		channel = b00001111
	}
	if pitch > b01111111 {
		pitch = b01111111
	}
	if velocity > b01111111 {
		velocity = b01111111
	}

	return []byte{
		(channel & b00001111) | b10010000,
		pitch,
		velocity,
	}
}

// NoteOff sends the signal to end playing a note on the given
// channel using the given pitch where velocity is usually 0.
func NoteOff(channel, pitch, velocity byte) []byte {
	if channel > b00001111 {
		channel = b00001111
	}
	if pitch > b01111111 {
		pitch = b01111111
	}
	if velocity > b01111111 {
		velocity = b01111111
	}

	return []byte{
		(channel & b00001111) | b10000000,
		pitch,
		velocity,
	}
}

// DrumOn sends the signal to start playing a note on
// channel 10 using the given drum type and velocity.
// Drum Notes do not have a pitch. The pitch value of
// NoteOn selects the kind of drum, the channel for drums
// is usually 9 (channel 10).
// It is essentially a specialized NoteOn.
func DrumOn(drum, velocity byte) []byte {
	if drum > b01111111 {
		drum = b01111111
	}
	if velocity > b01111111 {
		velocity = b01111111
	}

	return []byte{
		0x99, // Note on + Channel 10 (9)
		drum,
		velocity,
	}
}

// DrumOff sends the signal to stop playing a note on
// channel 10 using the given drum type where velocity
// is usually 0.
// It is essentially a specialized Noteff.
func DrumOff(drum, velocity byte) []byte {
	if drum > b01111111 {
		drum = b01111111
	}
	if velocity > b01111111 {
		velocity = b01111111
	}

	return []byte{
		0x89, // Note off + Channel 10 (9)
		drum & b01111111,
		velocity & b01111111,
	}
}

// Controller sends a controller message on the given
// channel with the given controller number and a value.
// Common controller numbers:
//     0 = Sound bank selection (MSB)
//     1 = Modulation wheel, often assigned to a vibrato or tremolo effect.
//     7 = Volume level of the instrument
//    10 = Panoramic (0 = left; 64 = center; 127 = right)
//    11 = Expression (sometimes used also for volume control or similar, depending on the synthesizer)
//    32 = Sound bank selection (LSB)
//    64 = Sustain pedal (0 = no pedal; >= 64 => pedal ON)
//   121 = All controllers off (this message clears all the controller values for this channel, back to their default values)
//   123 = All notes off (this message stops all the notes that are currently playing)
func Controller(channel, controller, value byte) []byte {
	if channel > b00001111 {
		channel = b00001111
	}
	if controller > b01111111 {
		controller = b01111111
	}
	if value > b01111111 {
		value = b01111111
	}

	return []byte{
		(channel & b00001111) | b10110000,
		controller & b01111111,
		value & b01111111,
	}
}

// ProgramChange sends a request to change the instrument to the
// one with the given value.
// See https://www.midi.org/specifications/item/gm-level-1-sound-set
// for the list of common instruments
func ProgramChange(channel, value byte) []byte {
	if channel > b00001111 {
		channel = b00001111
	}
	if value > b01111111 {
		value = b01111111
	}

	return []byte{
		(channel & b00001111) | b11000000,
		value,
	}
}

// BendPitch allows manipulating the pitch of the entire channel by the
// given 14bit value from 0x0 to 0x3fff with 0x2000 being the middle value
// meaning no change in pitch
func BendPitch(channel byte, value uint16) []byte {
	if channel > b00001111 {
		channel = b00001111
	}
	if value > b0011111111111111 {
		value = b0011111111111111
	}

	return []byte{
		(channel & b00001111) | b11100000,
		byte(value >> 7),
		byte(value & b01111111),
	}
}

// Reset sends a reset message which should make the synthesizer
// stop playing notes and go to its default settings.
func Reset() []byte {
	return []byte{b11111111}
}
