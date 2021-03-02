package device

// Thx to https://github.com/gvalkov/golang-evdev

import (
	"bytes"
	"encoding/binary"
	"os"
)

// InputDevice represents a Linux input device from which events can be read.
type InputDevice struct {
	Fn string // path to input device (devnode)

	Name string   // device name
	Phys string   // physical topology of device
	File *os.File // an open file handle to the input device

	Bustype uint16 // bus type identifier
	Vendor  uint16 // vendor identifier
	Product uint16 // product identifier
	Version uint16 // version identifier

	EvdevVersion int // evdev protocol version
}

// Open an evdev input device.
func Open(devnode string) (*InputDevice, error) {
	f, err := os.Open(devnode)
	if err != nil {
		return nil, err
	}

	dev := InputDevice{}
	dev.Fn = devnode
	dev.File = f

	return &dev, nil
}

// Read and return a slice of input events from device.
func (dev *InputDevice) Read() ([]InputEvent, error) {
	events := make([]InputEvent, 16)
	buffer := make([]byte, eventsize*16)

	_, err := dev.File.Read(buffer)
	if err != nil {
		return events, err
	}

	b := bytes.NewBuffer(buffer)
	err = binary.Read(b, binary.LittleEndian, &events)
	if err != nil {
		return events, err
	}

	// remove trailing structures
	for i := range events {
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}

	return events, err
}
