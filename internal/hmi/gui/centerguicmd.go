package gui

import (
	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func appendCenterGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	centerGroupGUICmd := &cobra.Command{
		Use:   "group",
		Short: "Center the selected group of windows",
		Long:  `TODO`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display group labels on windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(0)

			forceFloat, err := cmd.Flags().GetBool("force-float")
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the force-float flag")

				return
			}

			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group to snap")

				return
			}

			groupWindows := selectedGroup.GetWindows()
			for it := groupWindows.Iterator(); it.Next(); {
				curWin := it.Value().(*windowmgt.Window)

				if curErr := windowmgt.CenterWindow(curWin, forceFloat); curErr != nil {
					internal.NormalLogger.Warnf("Unable to center window '%s'", curWin.Name())
					err = eris.Wrapf(err, "Unable to center windows of group '%s'", selectedGroup.Name())
				}
			}

			return
		},
	}
	centerGroupGUICmd.Flags().Bool("force-float", false, "Make the window float 🤡")

	parentCmd.AddCommand(centerGroupGUICmd)
}

func AppendCenterGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	centerGUICmd := &cobra.Command{
		Use:   "center",
		Short: "Center elements",
	}

	appendCenterGroupGUICmd(centerGUICmd)

	parentCmd.AddCommand(centerGUICmd)
}
