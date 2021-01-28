package gui

import (
	"github.com/spf13/cobra"
)

// guiCmd represents the gui command.
func AppendGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	guiCmd := &cobra.Command{
		Use:   "gui",
		Short: "Commands to trigger the gui",
		Long: `A set of commands related to groups management.

		Se below the different available commands`,
	}

	AppendSnapGUICmd(guiCmd)
	AppendGroupMgtGUICmd(guiCmd)
	AppendVisibilityMgtGUICmd(guiCmd)
	AppendCenterGUICmd(guiCmd)

	parentCmd.AddCommand(guiCmd)
}
