package talkiepi

import (
	"crypto/tls"
	"github.com/dchote/gpio"
	"github.com/dchote/gumble/gumble"
	"github.com/dchote/gumble/gumbleopenal"
)


// Raspberry Pi GPIO pin assignments (CPU pin definitions)
const (
	OnlineLEDPin       uint = 18
	ParticipantsLEDPin uint = 23
	TransmitLEDPin     uint = 24
    TransmitLED2Pin    uint = 7		// not needed
	TransmitLED3Pin    uint = 12	// not needed
	ButtonPin          uint = 25
)

type Talkiepi struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   		string
	AltAddress   	string
	TLSConfig 		tls.Config

	ConnectAttempts uint

	Stream *gumbleopenal.Stream

	ChannelName    	string
	IsConnected    	bool
	IsTransmitting 	bool

	GPIOEnabled     bool
	OnlineLED       gpio.Pin
	ParticipantsLED gpio.Pin
	TransmitLED     gpio.Pin
	Transmit2LED    gpio.Pin
	Transmit3LED    gpio.Pin
	Button          gpio.Pin
	ButtonState     uint
}
