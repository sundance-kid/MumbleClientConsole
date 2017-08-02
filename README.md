# RasPi stage intercom

RasPi stage intercom is a wired and wireless stage intercom system.
It used a headless Raspberry Pi in a 3d printed enclosure as a belt pack,
communicates via wired Ethernet (ability to power the Pi via Power Over Ethernet)
and/or regular Wifi with a local mumble server.
The Mumble client is written in Go.
The project is based on the [talkiepi](https://github.com/dchote/talkiepi/), a walkie talkie style use off the Pi.
Both use GPIO pins for push to talk and LED status.
The Mumble client is a fork of [barnard](https://github.com/layeh/barnard),
which was a great jump off point for me to learn golang and have something that worked relatively quickly.


## 3D printable enclosure

In the stl directory are the stl files for an enclosure for the Raspberry Pi 3 Model B
You will also need a a random USB sound card (to have a microphone input) and a headset + power supply.
(PoE to micro USB extractors exist to only require the one Ethernet cable and no local power supply or power bank)

On the top side are 3 holes for 3mm LEDs in LED-fittings or 5mm LEDs without.
Next to it is a 2.5mm TRS connector for an external Push-To-Talk button. You can use an M6 thread-cutter in the 3D printed plastic to make it screw in perfectly.
The big hole is for a 19mm illuminated Push-To-Talk button.


## Installing

For an easy installation you can use the SD-card image provided.
Just change cmdline.txt for the IP-address,
start_intercom for the mumble-server and username
and motd.txt for a custom greeting uppon SSH login.
All these files are accessible on the FAT partition without access to a Linux machine that can read the main file system.
The image will mount all file systems read-only. So it is safe to turn off at any time.

You may want to change the password of the default Raspberry Pi user "pi" from "raspberry" to something else and create an individual SSH host key. 

For the manual installation onto an existing Raspberry Pi, there is an install guide [here](doc/README.md).


## GPIO and LEDs

You can edit your pin assignments in `talkiepi.go`
```go
const (
	OnlineLEDPin       uint = 18
	ParticipantsLEDPin uint = 23
	TransmitLEDPin     uint = 24
	ButtonPin          uint = 25
)
```

You can find the pinout for a Raspberry Pi 3 here: https://www.element14.com/community/servlet/JiveServlet/previewBody/73950-102-11-339300/pi3_gpio.png

Here is a basic schematic of how I am currently controlling the LEDs and pushbutton:

![schematic](doc/gpio_diagram.png "GPIO Diagram")

A red LED needs 1.8V, so we need 330Ω. 2x 634Ω in parallel are 317Ω. This is close enough. (317Ω * 0.005A = 1.585V; 3.3V - 1.585V = 1.715V)

A green LED needs 2.4V, so we need 220Ω. 2x 440Ω in parallel are obviously perfect.

A yellow LED needs 2.1V, so we need 270Ω. 2x 130Ω in series are close enough.

A blue or white LED needs 3.0V. This needs only a very small resistor (3.3V-3.0V=0.3V => X * 0.005A = 0.3V => X=0.3V/0.005A=60Ω).

Remember, the long side on an LED is +. That's the side that faces VCC=+3.3V on the GPIO pin. The short side faces GND with the resistor in between.

The Push-To-Talk button can also trigger an LED via a chain "3.3V---LED---Led resistor---gpio as input --- button --- GND"

So for a push button with embedded LED we connect:
(+) to 3.3V
(-) via LED-resistor to (N.O.=Normally Open) and to the GPIO pin
(C=Common) to GND


## License

MPL 2.0

## Author

- RasPi-Stage-Intercom - Marcus Wolschon (<Marcus@Wolschon.biz>)
- talkiepi - [Daniel Chote](https://github.com/dchote)
- Barnard,Gumble Author - Tim Cooper (<tim.cooper@layeh.com>)

