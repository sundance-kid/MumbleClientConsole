package talkiepi

import (
	"fmt"
	"time"
	"github.com/dchote/gpio"
	"github.com/stianeikeland/go-rpio"
)


func (b *Talkiepi) initGPIO() {
	// we need to pull in rpio to pullup our button pin
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		b.GPIOEnabled = false
		return
	} else {
		b.GPIOEnabled = true
	}

	ButtonPinPullUp := rpio.Pin(ButtonPin)
	ButtonPinPullUp.PullUp()
	rpio.Close()

	// Polling PTT button
	b.Button = gpio.NewInput(ButtonPin)		// Assign the PTT button GPIO-Pin
	
	go func() {														// start polling thread
		for {														// forever loop...
			currentState, err := b.Button.Read()
			if currentState != b.ButtonState && err == nil {
				b.ButtonState = currentState

				if b.Stream != nil {
					if b.ButtonState == 1 {
						fmt.Printf("PTT-Button is released...\n")
						b.TransmitStop()
					} else {
						fmt.Printf("PTT-Button is pressed...\n")
						b.TransmitStart()
					}
				}

			}
			time.Sleep(500 * time.Millisecond)
		} // for
	}()

	// then we can do our GPIO stuff
	b.OnlineLED = gpio.NewOutput(OnlineLEDPin, false)
	b.ParticipantsLED = gpio.NewOutput(ParticipantsLEDPin, false)
	b.TransmitLED = gpio.NewOutput(TransmitLEDPin, false)
	b.Transmit2LED = gpio.NewOutput(TransmitLEDPin, false)
	b.Transmit3LED = gpio.NewOutput(TransmitLEDPin, false)
}


func (b *Talkiepi) LEDOn(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}
	LED.High()
}


func (b *Talkiepi) LEDOff(LED gpio.Pin) {
	if b.GPIOEnabled == false {
		return
	}
	LED.Low()
}

func (b *Talkiepi) LEDOffAll() {
	if b.GPIOEnabled == false {
		return
	}
	b.LEDOff(b.OnlineLED)
	b.LEDOff(b.ParticipantsLED)
	b.LEDOff(b.TransmitLED)
}
