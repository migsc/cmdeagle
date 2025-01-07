package config

import (
	"cmdeagle/types"
	"fmt"

	"github.com/charmbracelet/log"
)

// CommandVisitor defines what operations can be performed when visiting a command
type CommandVisitor interface {
	Visit(cmd *types.CommandDefinition, parent *types.CommandDefinition, path []string) error
	// Build(config *CmdeagleConfig, cmd *types.CommandDefinition, parent *types.CommandDefinition, path []string) error
}

// WalkCommands traverses the command tree with a visitor
func WalkCommands(commands *[]types.CommandDefinition, parent *types.CommandDefinition, visitor CommandVisitor, path []string) error {
	for i := range *commands {
		log.Debug("Walking command", "name", (*commands)[i].Name)
		cmd := &(*commands)[i]
		currentPath := append(path, cmd.Name)

		if visitor == nil {
			return fmt.Errorf("Somehow visitor is nil. This should never happen. Please report this bug.")
		}

		if err := visitor.Visit(cmd, parent, currentPath); err != nil {
			return err
		}

		if len(cmd.Commands) > 0 {
			if err := WalkCommands(&cmd.Commands, cmd, visitor, currentPath); err != nil {
				return err
			}
		}
	}
	return nil
}
