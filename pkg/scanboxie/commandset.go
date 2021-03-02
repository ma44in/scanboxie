package scanboxie

import (
	"bytes"
	"fmt"
	"html/template"
	"os/exec"
)

// CommandSets ...
type CommandSets map[string]CommandSet

// CommandSet represents a set of commands executed
// one after the other
type CommandSet []Command

// Command represents an command
type Command struct {
	CommandString string
	commandSlice  []string
}

// ExecuteCommands templates commands and executes them
func (cs *CommandSet) ExecuteCommands(templateData interface{}) error {
	var templatedCommands []Command

	// TEMPLATE
	// command is something like "echo {.Dir}"
	for _, command := range *cs {
		fmt.Printf("Template command: %s\n", command)
		t := template.Must(template.New("").Parse(command.CommandString))

		var templatedCommandString bytes.Buffer
		if err := t.Execute(&templatedCommandString, templateData); err != nil {
			return err
		}

		var templatedCommand Command // := Command{CommandString: templatedCommandString.String()}
		templatedCommand.CommandString = templatedCommandString.String()
		templatedCommand.commandSlice = Split(templatedCommandString.String(), ' ')

		templatedCommands = append(templatedCommands, templatedCommand)
	}

	// EXECUTE
	for _, command := range templatedCommands {
		fmt.Printf("EXECUTE COMMAND: %s\n", command.commandSlice)

		cmd := exec.Command(command.commandSlice[0], command.commandSlice[1:]...)
		out, err := cmd.CombinedOutput()
		fmt.Printf("COMMAND OUTPUT:\n%s\n", string(out))

		if err != nil {
			fmt.Printf("COMMAND ERROR: %v\n", err)
			return err
		}
	}

	return nil
}

// func (cs *CommandSet) ExecuteCommands() error {

// 	return nil
// }
