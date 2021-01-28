package gui

import (
	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/hmi"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func appendSnapWindowGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapWindowGUICmd := &cobra.Command{
		Use:   "window",
		Short: "Snap the focused window to the selected cell",
		Long:  `TODO`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return hmi.BindGridDefinitionFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cols := viper.GetUint32("grid.cols")
			rows := viper.GetUint32("grid.rows")

			forceFloat, err := cmd.Flags().GetBool("force-float")
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the force-float argument")

				return
			}

			focusedWin, err := hmi.GetTargetedWindow(nil)
			if err != nil {
				err = eris.Wrap(err, "Unable to retrieve the focused window")

				return
			}

			internal.NormalLogger.Debugf("Window '%s' is focused", focusedWin.Name())

			selectedRow, selectedCol, err := AskUserToSelectCell(rows, cols)
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the selected cell")

				return
			}

			if err = windowmgt.SnapWindow(focusedWin, rows, cols, selectedRow, selectedCol, forceFloat); err != nil {
				err = eris.Wrap(err, "Unable to snap the window")

				return
			}

			return
		},
	}

	snapWindowGUICmd.Flags().Bool("force-float", false, "Make the window float 🤡")
	hmi.AddGridDefinitionFlags(snapWindowGUICmd)

	parentCmd.AddCommand(snapWindowGUICmd)
}

func appendSnapGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapGroupGUICmd := &cobra.Command{
		Use:   "group",
		Short: "Snap the selected windows group to the selected cell",
		Long:  `TODO`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return hmi.BindGridDefinitionFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display group labels on windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(0)
			cols := viper.GetUint32("grid.columns")
			rows := viper.GetUint32("grid.rows")

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

			selectedRow, selectedCol, err := AskUserToSelectCell(rows, cols)
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the coordinates of the selected cell")

				return
			}

			groupWindows := selectedGroup.GetWindows()
			for it := groupWindows.Iterator(); it.Next(); {
				curWin := it.Value().(*windowmgt.Window)

				if curErr := windowmgt.SnapWindow(curWin, viper.GetUint32("grid.rows"), viper.GetUint32("grid.cols"),
					selectedRow, selectedCol, forceFloat); curErr != nil {
					internal.NormalLogger.Warnf("Unable to snap window '%s'", curWin.Name())
					err = eris.Wrapf(err, "Unable to windows of group '%s'", selectedGroup.Name())
				}
			}

			return
		},
	}
	snapGroupGUICmd.Flags().Bool("force-float", false, "Make the window float 🤡")

	hmi.AddGridDefinitionFlags(snapGroupGUICmd)

	parentCmd.AddCommand(snapGroupGUICmd)
}

func AppendSnapGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapGUICmd := &cobra.Command{
		Use:   "snap",
		Short: "Snap elements the selected cell",
	}
	appendSnapWindowGUICmd(snapGUICmd)
	appendSnapGroupGUICmd(snapGUICmd)

	parentCmd.AddCommand(snapGUICmd)
}
