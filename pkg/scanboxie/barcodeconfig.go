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
	modified             bool
}

// NewBarcodeConfig create a new BarcodeConfig from given json file
func NewBarcodeConfig(path string, availableCommandSets *CommandSets) (*BarcodeConfig, error) {
	var barcodeConfig BarcodeConfig

	// Load BarcodeActions from file
	yamlContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlContent, &barcodeConfig)
	if err != nil {
		return nil, err
	}

	barcodeConfig.path = path
	barcodeConfig.availableCommandSets = availableCommandSets
	barcodeConfig.modified = false

	return &barcodeConfig, nil
}

func (bc *BarcodeConfig) Save() error {
	yamlContent, err := yaml.Marshal(bc)
	if err != nil {
		return fmt.Errorf("could not marshal barcodeconfig, err: %v", err)
	}

	err = ioutil.WriteFile(bc.path, yamlContent, 0644)
	if err != nil {
		return fmt.Errorf("could not save %s, err: %v", bc.path, err)
	}

	bc.modified = false
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
	bc.modified = true
	return nil
}

func (bc *BarcodeConfig) RemoveBarcodeAction(barcode string) error {
	for i, barcodeAction := range bc.BarcodeActions {
		if barcodeAction.Barcode == barcode {
			bc.BarcodeActions = append(bc.BarcodeActions[:i], bc.BarcodeActions[i+1:]...)
			bc.modified = true
			return nil
		}
	}

	return fmt.Errorf("could not remove barcode action with barcode %s, because it was not found", barcode)
}

func (bc *BarcodeConfig) MoveBarcodeAction(barcode string, newIndex int) error {
	if newIndex < 0 || newIndex >= len(bc.BarcodeActions) {
		return fmt.Errorf("could not move barcode %s, because new index %d is out of range", barcode, newIndex)
	}

	barcodeAction := bc.GetBarcodeAction(barcode)
	if barcodeAction == nil {
		return fmt.Errorf("barcode action with barcode %s not found", barcode)
	}

	err := bc.RemoveBarcodeAction(barcode)
	if err != nil {
		return err
	}

	bc.BarcodeActions = append(bc.BarcodeActions[:newIndex], append([]*BarcodeAction{barcodeAction}, bc.BarcodeActions[newIndex:]...)...)
	return nil
}

func (bc *BarcodeConfig) IsModified() bool {
	return bc.modified
}
