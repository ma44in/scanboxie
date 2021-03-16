package scanboxie

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"image/color"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/code128"
)

// BarcodeConfig stores a map of barcodes to directories
type BarcodeConfig struct {
	path              string
	BarcodeActions    []BarcodeAction
	barcodeActionsMap map[string]*BarcodeAction
}

// GetBarcodeAction returns a BarcodeAction based on the given barcode
func (bc *BarcodeConfig) GetBarcodeAction(barcode string) *BarcodeAction {
	return bc.barcodeActionsMap[barcode]
}

// BarcodeAction ...
type BarcodeAction struct {
	Barcode         string
	ActionKey       string
	BookletPageType string
	Values          []string
}

// GetBarcodeSvg returns Barcode Image as SVG
func (ba BarcodeAction) GetBarcodeSvg() template.HTML {
	// Thx to: https://github.com/boombuler/barcode/issues/57

	bc, err := code128.Encode(ba.Barcode)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)

	height := 100
	bounds := bc.Bounds()
	//canvas.Startview(bounds.Dx(), height, 0, 0, 100, 100)
	//canvas.Startpercent(bounds.Dx(), 100, "viewBox=\"0 0 100 100\"")

	// TODO heigth 50  problem ... still scales with aspect ratio
	canvas.Startraw("preserveAspectRatio=\"none\"", fmt.Sprintf("viewBox=\"0 0 %d 50\"", bounds.Dx()))

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if bc.At(x, bounds.Min.Y) == color.Black {
			start := x
			x++

			for x < bounds.Max.X && bc.At(x, bounds.Min.Y) == color.Black {
				x++
			}

			canvas.Rect(start, 0, x-start, height, "fill:black")
		}
	}

	canvas.End()
	return template.HTML(buf.String())
}

// NewBarcodeConfig create a new BarcodeConfig from given json file
func NewBarcodeConfig(path string, watchForChanges bool) (*BarcodeConfig, error) {
	// barcodeConfig := BarcodeConfig{
	// 	path:           path,
	// 	BarcodeActions: make(map[string]BarcodeAction),
	// }

	var barcodeConfig BarcodeConfig
	barcodeConfig.path = path
	barcodeConfig.BarcodeActions = []BarcodeAction{}
	barcodeConfig.barcodeActionsMap = make(map[string]*BarcodeAction)

	barcodeConfig.readConfigfile()
	if watchForChanges {
		err := barcodeConfig.watchForChanges()
		if err != nil {
			fmt.Printf("could not watch changes. error: %v\n", err)
		}
	}

	return &barcodeConfig, nil
}

func (bc *BarcodeConfig) readConfigfile() error {
	barcodeActions, err := ReadCsv(bc.path)
	if err != nil {
		panic(err)
	}

	bc.BarcodeActions = barcodeActions

	// Add BarcodeAction to Map for easier lookup
	for _, barcodeAction := range barcodeActions {
		bc.barcodeActionsMap[barcodeAction.Barcode] = &barcodeAction
	}

	return nil
}

func (bc *BarcodeConfig) watchForChanges() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add(bc.path); err != nil {
		return err
	}

	fmt.Printf("watch for changes on %s", bc.path)

	go func() {
		for {
			select {
			case <-watcher.Events:
				fmt.Printf("%s changed", bc.path)
				bc.readConfigfile()
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Println("watcher error:", err)
				}
				//case <-bc.quit:
				//		watcher.Close()
				//	return
			}
		}
	}()
	return nil

	// // TODO does not work

	// // creates a new file watcher
	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	fmt.Println("ERROR", err)
	// }
	// defer watcher.Close()

	// done := make(chan bool)
	// go func() {
	// 	for {
	// 		select {
	// 		case event, ok := <-watcher.Events:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("event:", event)
	// 			if event.Op&fsnotify.Write == fsnotify.Write {
	// 				log.Println("modified file:", event.Name)
	// 			}
	// 		case err, ok := <-watcher.Errors:
	// 			if !ok {
	// 				return
	// 			}
	// 			log.Println("error:", err)
	// 		}
	// 	}
	// }()

	// fmt.Printf("Watch changes on %s", bc.path)
	// err = watcher.Add(bc.path)
	// if err != nil {
	// 	return fmt.Errorf("could not add file %s to watcher", bc.path)
	// }
	// <-done

	// return nil
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([]BarcodeAction, error) {

	// Open CSV file
	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	// Setup CSV Reader
	csvreader := csv.NewReader(csvfile)
	csvreader.Comment = '#'
	csvreader.FieldsPerRecord = -1 // -1 Variable Records

	// Read Line by Line
	barcodeActions := []BarcodeAction{}

	for {
		line, err := csvreader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("skip line due to error. line: %s, error: %v", line, err)
			continue
		}

		if len(line) < 3 {
			fmt.Printf("skip line due to missing fields. need %d fields got %d, line: %v\n", 4, len(line), line)
			continue
		}

		barcode := line[0]
		actionKey := line[1]
		bookletPageType := line[2]
		values := line[3:]

		barcodeAction := BarcodeAction{
			Barcode:         barcode,
			ActionKey:       actionKey,
			BookletPageType: bookletPageType,
			Values:          values,
		}

		barcodeActions = append(barcodeActions, barcodeAction)
	}

	return barcodeActions, nil
}
