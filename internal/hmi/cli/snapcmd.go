package cli

import (
	"strconv"

	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/hmi"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func appendSnapWindowCmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapWindowCmd := &cobra.Command{
		Args:  cobra.RangeArgs(2, 3), //nolint:gomnd //Obvious number of positional arguments
		Use:   "window <SELECTEDROW> <SELECTEDCOL> [XWINDOWID]",
		Short: "Snap a window to a cell of a grid",
		Long: `Move and resize a window to the given cell of a grid.
	
  If no X Window ID is given, it will snap the focused window.
  By default the command does nothing to non-floating window. 
	Use the force-float argument to make the window float before snapping it.

  The grid is defined by the 'cols' and 'rows' arguments. The position on which the window will be
  snapped is given by the 'SELECTEDROW' and 'SELECTEDCOL' arguments. 
	The size of the window is based on the resolution of the focused screen and depends on 
	the 'cols' and 'rows' arguments.`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return hmi.BindGridDefinitionFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedRowArg := args[0]
			selectedColArg := args[1]
			selectedRow, err := strconv.Atoi(selectedRowArg)
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the selected row '%s'", selectedRowArg)

				return
			}

			selectedCol, err := strconv.Atoi(selectedColArg)
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the selected column '%s'", selectedColArg)

				return
			}

			forceFloat, err := cmd.Flags().GetBool("force-float")
			if err != nil {
				err = eris.Wrap(internal.InvalidArgumentError, "Unable to retrieve the force-float flag")

				return
			}

			var xWinIDStr *string
			if len(args) == 3 { //nolint:gomnd //Obvious number of positional arguments
				xWinIDStr = &args[2]
			}

			var wm *windowmgt.WindowManager
			wm, err = windowmgt.GetWindowManagerInstance()
			if err != nil {
				err = eris.Wrapf(err, "Unable to snap the windows: cannot retrieved windowmanager")

				return
			}

			filterStrs, err := cmd.Flags().GetStringArray("filter")
			if err != nil {
				// TODO: deal with error
				return
			}
			windowsToProcess, err := GetWindowsToProcess(xWinIDStr, wm.GetWindows(func(w *windowmgt.Window) bool { return true }), filterStrs)

			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
				curWin := winIt.Value().(*windowmgt.Window)
				if err = windowmgt.SnapWindow(curWin, viper.GetUint32("grid.rows"),
					viper.GetUint32("grid.cols"), uint32(selectedRow), uint32(selectedCol), forceFloat); err != nil {
					err = eris.Wrapf(err, "Unable to snap the window '%s'", curWin.Name())

					return
				}
			}

			return
		},
	}

	snapWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")
	snapWindowCmd.Flags().Bool("force-float", false, "Make the window float 🤡")
	hmi.AddGridDefinitionFlags(snapWindowCmd)

	parentCmd.AddCommand(snapWindowCmd)
}

func appendSnapGroupCmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapGroupCmd := &cobra.Command{
		Args:  cobra.ExactArgs(3), //nolint:gomnd //Obvious number of positional arguments
		Use:   "group <SELECTEDROW> <SELECTEDCOL> <GROUPID>",
		Short: "Snap a group of windows window to a cell of a grid",
		Long: `Move and resize a group of window to the given cell of a grid.
	
  By default the command does nothing to non-floating window. 
	Use the force-float argument to make the window float before snapping it.

  The grid is defined by the 'cols' and 'rows' arguments. The position on which thes windows will be
  snapped is given by the 'SELECTEDROW' and 'SELECTEDCOL' arguments. 
	The size of the window is based on the resolution of the focused screen and depends on 
	the 'cols' and 'rows' arguments.`,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return hmi.BindGridDefinitionFlags(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedRowArg := args[0]
			selectedColArg := args[1]

			selectedRow, err := strconv.Atoi(selectedRowArg)
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the selected row '%s'", selectedRowArg)

				return
			}

			selectedCol, err := strconv.Atoi(selectedColArg)
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the selected column '%s'", selectedColArg)

				return
			}

			groupID, err := strconv.Atoi(args[2])
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group id '%d'", groupID)

				return
			}

			forceFloat, err := cmd.Flags().GetBool("force-float")
			if err != nil {
				err = eris.Wrap(internal.InvalidArgumentError, "Unable to retrieve the force-float flag")

				return
			}

			wgm := windowmgt.GetWindowGroupManagerInstance()

			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
			if err != nil {
				err = eris.Wrapf(err, "Unable to retrieve the group with ID '%d'", groupID)

				return
			}

			filterStrs, err := cmd.Flags().GetStringArray("filter")
			windowsToProcess := group.GetWindows()
			if len(filterStrs) != 0 {
				windowsToProcess, err = FilterWindowsToProcess(windowsToProcess, filterStrs)
				if err != nil {
					err = eris.Wrapf(err, "Unable to snap the windows: cannot filter windows")

					return
				}
			}

			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
				curWin := winIt.Value().(*windowmgt.Window)
				err = windowmgt.SnapWindow(curWin, viper.GetUint32("grid.rows"),
					viper.GetUint32("grid.cols"), uint32(selectedRow), uint32(selectedCol), forceFloat)
				if err != nil {
					err = eris.Wrapf(err, "Unable to snap the window '%s'", curWin.Name())

					return
				}
			}

			return
		},
	}

	snapGroupCmd.Flags().Bool("force-float", false, "Make the windows float 🤡")
	snapGroupCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")
	hmi.AddGridDefinitionFlags(snapGroupCmd)

	parentCmd.AddCommand(snapGroupCmd)
}

func AppendSnapCmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapCmd := &cobra.Command{
		Use:   "snap",
		Short: "Snap elements to the given cell",
	}

	appendSnapWindowCmd(snapCmd)
	appendSnapGroupCmd(snapCmd)

	parentCmd.AddCommand(snapCmd)
}
