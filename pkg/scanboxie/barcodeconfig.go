package scanboxie

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
)

// BarcodeConfig stores a map of barcodes to directories
type BarcodeConfig struct {
	path           string
	BarcodeActions map[string]BarcodeAction
}

// BarcodeAction ...
type BarcodeAction struct {
	ActionKey string
	Values    []string
}

// NewBarcodeConfig create a new BarcodeConfig from given json file
func NewBarcodeConfig(path string, watchForChanges bool) (*BarcodeConfig, error) {
	barcodeConfig := BarcodeConfig{
		path:           path,
		BarcodeActions: make(map[string]BarcodeAction),
	}

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
func ReadCsv(filename string) (map[string]BarcodeAction, error) {

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
	barcodeActions := make(map[string]BarcodeAction)

	for {
		line, err := csvreader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("skip line due to error. line: %s, error: %v", line, err)
			continue
		}

		if len(line) < 2 {
			fmt.Printf("skip line due to missing fields. need %d fields got %d, line: %v", 2, len(line), line)
			continue
		}

		barcode := line[0]
		actionKey := line[1]
		values := line[2:]

		barcodeAction := BarcodeAction{
			ActionKey: actionKey,
			Values:    values,
		}

		barcodeActions[barcode] = barcodeAction
	}

	return barcodeActions, nil
}
