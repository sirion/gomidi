
# gomidi - midi tools written in go

gomidi is a go module that contains the following packages/programs:


Path | Package
---- | ----
app/cmd/miDiMacro | Turns a midi- into a macro-keyboard with configurable key combinations assigned to midi notes and controller keys
lib/midi | package containing helper functions to create midi byte slices (midi commands)
lib/launchpadmini | package containing the LaunchpadMini struct which contains functions to read from the midi controller and set its state (turn colored button lights on and off, send text)

### miDiMacro

Turn your midi keyboard/controller into a macro keyboard.

The program takes its configuration from ```~/.config.midimacro/config.json```

Take a look at the [example configuration](/sirion/gomidi/blob/master/app/cmd/miDiMacro/config-example/config.json).

_(Only tested with the Launchpad Mini)_

### midi

The ```midi``` package only contains helper functions that create midi-messages (byte-slices). It is essentially the code version of what I learned from reading http://www.music-software-development.com/midi-tutorial.html.

Read the documentation at https://godoc.org/github.com/sirion/gomidi/lib/midi.

### launchpadmini

The ```launchpadmini``` package contains the ```LaunchpadMini``` struct which can be created by calling ```launchpadmini.New(devicePath)``` to read key presses from the midi keyboard and control the button lights.

Read the documentation at https://godoc.org/github.com/sirion/gomidi/lib/launchpadmini.

## About:

This project started when I realized that I am never going to use my [Launchpad Mini](https://amzn.to/2SdAHys)* for its intended purpose so I decided to write a program to turn it into a macro keyboard.

In order to keep me motivated I developed the program on my [live stream](https://abovethelawn.de).


---
_*) Amazon Partner Link_
