package launchpadmini

// Colors supported by the LaunchPad Mini
const (
	ColorOff byte = 12

	ColorRedLow      byte = 13
	ColorRedFull     byte = 15
	ColorRedFlashing byte = 11

	ColorGreenLow      byte = 28
	ColorGreenFull     byte = 60
	ColorGreenFlashing byte = 56

	ColorAmberLow      byte = 29
	ColorAmberFull     byte = 63
	ColorAmberFlashing byte = 59

	ColorYellowFull     byte = 62
	ColorYellowFlashing byte = 58
)

// These constants describe all buttons on the Launchpad Mini.
// The constants represent the hardware/midi byte values for the (Grid) Buttons.
// The Live-Buttons are mapped to values +100 because of an overlap in the byte value of LiveButton1 with ButtonG
const (
	LiveButton1 byte = 104 + 100
	LiveButton2 byte = 105 + 100
	LiveButton3 byte = 106 + 100
	LiveButton4 byte = 107 + 100
	LiveButton5 byte = 108 + 100
	LiveButton6 byte = 109 + 100
	LiveButton7 byte = 110 + 100
	LiveButton8 byte = 111 + 100

	ButtonA1 byte = 0
	ButtonA2 byte = 1
	ButtonA3 byte = 2
	ButtonA4 byte = 3
	ButtonA5 byte = 4
	ButtonA6 byte = 5
	ButtonA7 byte = 6
	ButtonA8 byte = 7
	ButtonA  byte = 8

	ButtonB1 byte = 16
	ButtonB2 byte = 17
	ButtonB3 byte = 18
	ButtonB4 byte = 19
	ButtonB5 byte = 20
	ButtonB6 byte = 21
	ButtonB7 byte = 22
	ButtonB8 byte = 23
	ButtonB  byte = 24

	ButtonC1 byte = 32
	ButtonC2 byte = 33
	ButtonC3 byte = 34
	ButtonC4 byte = 35
	ButtonC5 byte = 36
	ButtonC6 byte = 37
	ButtonC7 byte = 38
	ButtonC8 byte = 39
	ButtonC  byte = 40

	ButtonD1 byte = 48
	ButtonD2 byte = 49
	ButtonD3 byte = 50
	ButtonD4 byte = 51
	ButtonD5 byte = 52
	ButtonD6 byte = 53
	ButtonD7 byte = 54
	ButtonD8 byte = 55
	ButtonD  byte = 56

	ButtonE1 byte = 64
	ButtonE2 byte = 65
	ButtonE3 byte = 66
	ButtonE4 byte = 67
	ButtonE5 byte = 68
	ButtonE6 byte = 69
	ButtonE7 byte = 70
	ButtonE8 byte = 71
	ButtonE  byte = 72

	ButtonF1 byte = 80
	ButtonF2 byte = 81
	ButtonF3 byte = 82
	ButtonF4 byte = 83
	ButtonF5 byte = 84
	ButtonF6 byte = 85
	ButtonF7 byte = 86
	ButtonF8 byte = 87
	ButtonF  byte = 88

	ButtonG1 byte = 96
	ButtonG2 byte = 97
	ButtonG3 byte = 98
	ButtonG4 byte = 99
	ButtonG5 byte = 100
	ButtonG6 byte = 101
	ButtonG7 byte = 102
	ButtonG8 byte = 103
	ButtonG  byte = 104

	ButtonH1 byte = 112
	ButtonH2 byte = 113
	ButtonH3 byte = 114
	ButtonH4 byte = 115
	ButtonH5 byte = 116
	ButtonH6 byte = 117
	ButtonH7 byte = 118
	ButtonH8 byte = 119
	ButtonH  byte = 120
)

/**
 * Pseudo constants
 */

// textCmdPrefix and textCmdSuffix are the byte sequences that need to be sent to the launchpad to create a text output
var textCmdPrefix = []byte{240, 0, 32, 41, 9}
var textCmdSuffix = []byte{247}

// ButtonNames is a pseudo-constant map to convert button name strings to the actual byte value
var ButtonValues = map[string]byte{
	"LiveButton1": LiveButton1,
	"LiveButton2": LiveButton2,
	"LiveButton3": LiveButton3,
	"LiveButton4": LiveButton4,
	"LiveButton5": LiveButton5,
	"LiveButton6": LiveButton6,
	"LiveButton7": LiveButton7,
	"LiveButton8": LiveButton8,

	"ButtonA1": ButtonA1,
	"ButtonA2": ButtonA2,
	"ButtonA3": ButtonA3,
	"ButtonA4": ButtonA4,
	"ButtonA5": ButtonA5,
	"ButtonA6": ButtonA6,
	"ButtonA7": ButtonA7,
	"ButtonA8": ButtonA8,
	"ButtonA":  ButtonA,

	"ButtonB1": ButtonB1,
	"ButtonB2": ButtonB2,
	"ButtonB3": ButtonB3,
	"ButtonB4": ButtonB4,
	"ButtonB5": ButtonB5,
	"ButtonB6": ButtonB6,
	"ButtonB7": ButtonB7,
	"ButtonB8": ButtonB8,
	"ButtonB":  ButtonB,

	"ButtonC1": ButtonC1,
	"ButtonC2": ButtonC2,
	"ButtonC3": ButtonC3,
	"ButtonC4": ButtonC4,
	"ButtonC5": ButtonC5,
	"ButtonC6": ButtonC6,
	"ButtonC7": ButtonC7,
	"ButtonC8": ButtonC8,
	"ButtonC":  ButtonC,

	"ButtonD1": ButtonD1,
	"ButtonD2": ButtonD2,
	"ButtonD3": ButtonD3,
	"ButtonD4": ButtonD4,
	"ButtonD5": ButtonD5,
	"ButtonD6": ButtonD6,
	"ButtonD7": ButtonD7,
	"ButtonD8": ButtonD8,
	"ButtonD":  ButtonD,

	"ButtonE1": ButtonE1,
	"ButtonE2": ButtonE2,
	"ButtonE3": ButtonE3,
	"ButtonE4": ButtonE4,
	"ButtonE5": ButtonE5,
	"ButtonE6": ButtonE6,
	"ButtonE7": ButtonE7,
	"ButtonE8": ButtonE8,
	"ButtonE":  ButtonE,

	"ButtonF1": ButtonF1,
	"ButtonF2": ButtonF2,
	"ButtonF3": ButtonF3,
	"ButtonF4": ButtonF4,
	"ButtonF5": ButtonF5,
	"ButtonF6": ButtonF6,
	"ButtonF7": ButtonF7,
	"ButtonF8": ButtonF8,
	"ButtonF":  ButtonF,

	"ButtonG1": ButtonG1,
	"ButtonG2": ButtonG2,
	"ButtonG3": ButtonG3,
	"ButtonG4": ButtonG4,
	"ButtonG5": ButtonG5,
	"ButtonG6": ButtonG6,
	"ButtonG7": ButtonG7,
	"ButtonG8": ButtonG8,
	"ButtonG":  ButtonG,

	"ButtonH1": ButtonH1,
	"ButtonH2": ButtonH2,
	"ButtonH3": ButtonH3,
	"ButtonH4": ButtonH4,
	"ButtonH5": ButtonH5,
	"ButtonH6": ButtonH6,
	"ButtonH7": ButtonH7,
	"ButtonH8": ButtonH8,
	"ButtonH":  ButtonH,
}

var ButtonNames = map[byte]string{
	LiveButton1: "LiveButton1",
	LiveButton2: "LiveButton2",
	LiveButton3: "LiveButton3",
	LiveButton4: "LiveButton4",
	LiveButton5: "LiveButton5",
	LiveButton6: "LiveButton6",
	LiveButton7: "LiveButton7",
	LiveButton8: "LiveButton8",

	ButtonA1: "ButtonA1",
	ButtonA2: "ButtonA2",
	ButtonA3: "ButtonA3",
	ButtonA4: "ButtonA4",
	ButtonA5: "ButtonA5",
	ButtonA6: "ButtonA6",
	ButtonA7: "ButtonA7",
	ButtonA8: "ButtonA8",
	ButtonA:  " ButtonA",

	ButtonB1: "ButtonB1",
	ButtonB2: "ButtonB2",
	ButtonB3: "ButtonB3",
	ButtonB4: "ButtonB4",
	ButtonB5: "ButtonB5",
	ButtonB6: "ButtonB6",
	ButtonB7: "ButtonB7",
	ButtonB8: "ButtonB8",
	ButtonB:  " ButtonB",

	ButtonC1: "ButtonC1",
	ButtonC2: "ButtonC2",
	ButtonC3: "ButtonC3",
	ButtonC4: "ButtonC4",
	ButtonC5: "ButtonC5",
	ButtonC6: "ButtonC6",
	ButtonC7: "ButtonC7",
	ButtonC8: "ButtonC8",
	ButtonC:  " ButtonC",

	ButtonD1: "ButtonD1",
	ButtonD2: "ButtonD2",
	ButtonD3: "ButtonD3",
	ButtonD4: "ButtonD4",
	ButtonD5: "ButtonD5",
	ButtonD6: "ButtonD6",
	ButtonD7: "ButtonD7",
	ButtonD8: "ButtonD8",
	ButtonD:  " ButtonD",

	ButtonE1: "ButtonE1",
	ButtonE2: "ButtonE2",
	ButtonE3: "ButtonE3",
	ButtonE4: "ButtonE4",
	ButtonE5: "ButtonE5",
	ButtonE6: "ButtonE6",
	ButtonE7: "ButtonE7",
	ButtonE8: "ButtonE8",
	ButtonE:  " ButtonE",

	ButtonF1: "ButtonF1",
	ButtonF2: "ButtonF2",
	ButtonF3: "ButtonF3",
	ButtonF4: "ButtonF4",
	ButtonF5: "ButtonF5",
	ButtonF6: "ButtonF6",
	ButtonF7: "ButtonF7",
	ButtonF8: "ButtonF8",
	ButtonF:  " ButtonF",

	ButtonG1: "ButtonG1",
	ButtonG2: "ButtonG2",
	ButtonG3: "ButtonG3",
	ButtonG4: "ButtonG4",
	ButtonG5: "ButtonG5",
	ButtonG6: "ButtonG6",
	ButtonG7: "ButtonG7",
	ButtonG8: "ButtonG8",
	ButtonG:  " ButtonG",

	ButtonH1: "ButtonH1",
	ButtonH2: "ButtonH2",
	ButtonH3: "ButtonH3",
	ButtonH4: "ButtonH4",
	ButtonH5: "ButtonH5",
	ButtonH6: "ButtonH6",
	ButtonH7: "ButtonH7",
	ButtonH8: "ButtonH8",
	ButtonH:  " ButtonH",
}

// ColorNames is a pseudo-constant map to convert color name strings to the actual byte value
var ColorNames = map[string]byte{
	"ColorOff": ColorOff,

	"ColorRedLow":      ColorRedLow,
	"ColorRedFull":     ColorRedFull,
	"ColorRedFlashing": ColorRedFlashing,

	"ColorGreenLow":      ColorGreenLow,
	"ColorGreenFull":     ColorGreenFull,
	"ColorGreenFlashing": ColorGreenFlashing,

	"ColorAmberLow":      ColorAmberLow,
	"ColorAmberFull":     ColorAmberFull,
	"ColorAmberFlashing": ColorAmberFlashing,

	"ColorYellowFull":     ColorYellowFull,
	"ColorYellowFlashing": ColorYellowFlashing,
}

var ColorValues = map[byte]string{
	ColorOff: "ColorOff",

	ColorRedLow:      "ColorRedLow",
	ColorRedFull:     "ColorRedFull",
	ColorRedFlashing: "ColorRedFlashing",

	ColorGreenLow:      "ColorGreenLow",
	ColorGreenFull:     "ColorGreenFull",
	ColorGreenFlashing: "ColorGreenFlashing",

	ColorAmberLow:      "ColorAmberLow",
	ColorAmberFull:     "ColorAmberFull",
	ColorAmberFlashing: "ColorAmberFlashing",

	ColorYellowFull:     "ColorYellowFull",
	ColorYellowFlashing: "ColorYellowFlashing",
}
