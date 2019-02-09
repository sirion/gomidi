package main

import (
	"fmt"
	"time"

	l "github.com/sirion/gomidi/lib/launchpadmini"
)

func main() {
	lp := l.New("auto")

	running := true
	in := lp.Listen()

	go func() {
		for running {
			button := <-in
			fmt.Printf("Button %s pressed\n", l.ButtonNames[button])

			if button == l.ButtonH {
				running = false
			}
		}
	}()

	/*
		lp.Text("Start", l.ColorAmberFull)

		time.Sleep(5 * time.Second)

		lp.Button(l.ButtonA1, l.ColorGreenFull)
		lp.Button(l.ButtonA2, l.ColorAmberFull)
		lp.Button(l.ButtonA3, l.ColorRedFull)

		lp.Grid(1, 0, l.ColorGreenFull)
		lp.Grid(1, 1, l.ColorAmberFull)
		lp.Grid(1, 2, l.ColorRedFull)

		lp.Live(0, l.ColorGreenFull)
		lp.Live(1, l.ColorAmberFull)
		lp.Live(2, l.ColorRedFull)

		lp.Button(l.LiveButton5, l.ColorGreenFull)
		lp.Button(l.LiveButton6, l.ColorAmberFull)
		lp.Button(l.LiveButton7, l.ColorRedFull)

		time.Sleep(5 * time.Second)

		lp.Reset()

		time.Sleep(time.Second)

		lp.AllOn(125)
		time.Sleep(500 * time.Millisecond)
		lp.AllOn(126)
		time.Sleep(500 * time.Millisecond)
		lp.AllOn(127)
		time.Sleep(500 * time.Millisecond)
		lp.AllOn(126)
		time.Sleep(500 * time.Millisecond)
		lp.AllOn(125)

		for i := 0; i < 20; i++ {
			lp.AllOn(125)
			time.Sleep(150 * time.Millisecond)
			lp.AllOn(127)
			time.Sleep(150 * time.Millisecond)
		}

		lp.Button(l.ButtonA, l.ColorGreenFull)
		lp.Flashing(true)

		time.Sleep(5 * time.Second)

		buttons := make(map[byte]byte, 10)
		buttons[l.ButtonA1] = l.ColorGreenFlashing
		buttons[l.ButtonB2] = l.ColorAmberFlashing
		buttons[l.ButtonC3] = l.ColorRedFlashing
		buttons[l.ButtonD4] = l.ColorGreenFlashing
		buttons[l.ButtonE5] = l.ColorAmberFlashing
		buttons[l.ButtonF6] = l.ColorRedFlashing

		lp.RapidUpdate(buttons)

		time.Sleep(5 * time.Second)
	*/
	lp.Reset()

	buttons1 := make(map[byte]byte, 6)
	buttons2 := make(map[byte]byte, 6)

	buttons1[l.ButtonA5] = l.ColorGreenFull
	buttons1[l.ButtonB5] = l.ColorAmberFull
	buttons1[l.ButtonC5] = l.ColorRedFull
	buttons1[l.ButtonA6] = l.ColorOff
	buttons1[l.ButtonB6] = l.ColorOff
	buttons1[l.ButtonC6] = l.ColorOff

	buttons2[l.ButtonA5] = l.ColorOff
	buttons2[l.ButtonB5] = l.ColorOff
	buttons2[l.ButtonC5] = l.ColorOff
	buttons2[l.ButtonA6] = l.ColorGreenFull
	buttons2[l.ButtonB6] = l.ColorAmberFull
	buttons2[l.ButtonC6] = l.ColorRedFull

	lp.BufferMode(l.BufferMode0)
	lp.RapidUpdate(buttons1)
	time.Sleep(time.Second)

	lp.BufferMode(l.BufferMode1)
	lp.RapidUpdate(buttons2)
	time.Sleep(time.Second)

	lp.BufferMode(l.BufferMode0)
	lp.RapidUpdate(buttons1)
	time.Sleep(time.Second)
	lp.BufferMode(l.BufferMode1)
	lp.RapidUpdate(buttons2)
	time.Sleep(time.Second)
	lp.BufferMode(l.BufferMode0)
	lp.RapidUpdate(buttons1)
	time.Sleep(time.Second)
	lp.BufferMode(l.BufferMode1)
	lp.RapidUpdate(buttons2)
	time.Sleep(time.Second)

	lp.BufferMode(l.BufferModeDefault)

	time.Sleep(5 * time.Second)

	lp.Reset()

	for running {
		time.Sleep(time.Second)
	}
}
