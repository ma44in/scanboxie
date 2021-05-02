package scanboxie

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// BarcodeAction ...
type BarcodeAction struct {
	Barcode    string `yaml:"Barcode"`
	Commandset string `yaml:"Commandset"`
	Value      string `yaml:"Value"`
	Pagetype   string `yaml:"Pagetype"`
}

// BarcodeConfig stores a map of barcodes to directories
type BarcodeConfig struct {
	path                 string
	BarcodeActions       []*BarcodeAction `yaml:"BarcodeActions"`
	availableCommandSets *CommandSets
}

// NewBarcodeConfig create a new BarcodeConfig from given json file
func NewBarcodeConfig(path string, availableCommandSets *CommandSets) (*BarcodeConfig, error) {
	var barcodeConfig BarcodeConfig

	// Load BarcodeActions from file
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &barcodeConfig)
	if err != nil {
		return nil, err
	}

	barcodeConfig.path = path
	barcodeConfig.availableCommandSets = availableCommandSets

	return &barcodeConfig, nil
}

func (bc *BarcodeConfig) AddBarcodeAction(barcode string, commandSetKey string, value string, pagetype string) error {
	barcodeAction := BarcodeAction{
		Barcode:    barcode,
		Commandset: commandSetKey,
		Value:      value,
		Pagetype:   pagetype,
	}

	commandSet := (*bc.availableCommandSets)[commandSetKey]
	if commandSet == nil {
		return fmt.Errorf("commandSet with key %s not found", commandSetKey)
	}

	if bc.GetBarcodeAction(barcode) != nil {
		return fmt.Errorf("barcode %s already exists in barcodeconfig", barcode)
	}

	bc.BarcodeActions = append(bc.BarcodeActions, &barcodeAction)
	return nil
}

// GetBarcodeAction returns a BarcodeAction based on the given barcode
func (bc *BarcodeConfig) GetBarcodeAction(barcode string) *BarcodeAction {
	for _, barcodeAction := range bc.BarcodeActions {
		if barcodeAction.Barcode == barcode {
			return barcodeAction
		}
	}

	return nil
}
