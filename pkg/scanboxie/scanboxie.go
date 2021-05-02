package scanboxie

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"scanboxie/pkg/device"
)

// Scanboxie app
type Scanboxie struct {
	BarcodeConfig *BarcodeConfig
	CommandSets   *CommandSets
}

// NewScanboxie returns a new Scanboxie app
func NewScanboxie(barcodeDirMapFilepath string, commandSets *CommandSets) (*Scanboxie, error) {
	barcodeConfig, err := NewBarcodeConfig(barcodeDirMapFilepath, commandSets)
	if err != nil {
		return nil, fmt.Errorf("could not load barcode config from %v. error: %v", barcodeConfig, err)
	}
	log.Printf("Loaded Barcode Config:\n%v\n\n", barcodeConfig)

	scanboxie := &Scanboxie{
		BarcodeConfig: barcodeConfig,
		CommandSets:   commandSets,
	}

	return scanboxie, nil
}

func (sb *Scanboxie) readEventsToPipe(targetDevice *device.InputDevice) *io.PipeReader {
	pr, pw := io.Pipe()

	var events []device.InputEvent
	var err error

	go func() {
		for {
			events, err = targetDevice.Read()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			shiftPressed := false

			// If a barcode is scanned, each character get received as group of events.
			// The group contains basically a Key-Down, Key-Up and End-Event.
			// A character in upper case has an aditionally Leftshift-Keypress event.
			for i := range events {
				myevent := &events[i]

				// Only look at KEY events with value = 1 (pressed)
				if myevent.Type == device.EV_KEY && myevent.Value == 1 {
					if myevent.Code == device.KEY_LEFTSHIFT {
						shiftPressed = true
					} else {
						str := getCharFromKeyEvent(myevent, shiftPressed)
						pw.Write([]byte(str))

						shiftPressed = false
					}
				}
			}
		}
	}()

	return pr
}

func (sb *Scanboxie) processEvents(pr *io.PipeReader, barcodeConfig *BarcodeConfig, commandSets *CommandSets) error {
	scanner := bufio.NewScanner(pr)
	for scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("Got input line: %s\n", input)

		// Lookup BarcodeAction for received barcode input
		barcodeAction := (*barcodeConfig).GetBarcodeAction(input)
		if barcodeAction != nil {
			templateData := struct {
				Value string
			}{
				Value: barcodeAction.Value,
			}

			commandSet := (*commandSets)[barcodeAction.Commandset]

			err := commandSet.ExecuteCommands(templateData)
			if err != nil {
				fmt.Printf("error execution command set with key %s. error: %v", barcodeAction.Commandset, err)
			}
		}
	}

	return nil
}

// ListenAndProcessEvents listens for input events on the given
// input path (e. g. /dev/input/event0)
func (sb *Scanboxie) ListenAndProcessEvents(eventPath string) error {

	// Open input device
	fmt.Println("Open device")
	targetDevice, _ := device.Open(eventPath)
	if targetDevice == nil {
		return fmt.Errorf("target device not found")
	}
	fmt.Printf("Found device: %s\n", targetDevice.Name)

	pr := sb.readEventsToPipe(targetDevice)
	err := sb.processEvents(pr, sb.BarcodeConfig, sb.CommandSets)

	return err
}
