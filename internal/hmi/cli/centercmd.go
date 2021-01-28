package cli

import (
	"strconv"

	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/hmi"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func appendCenterWindowCmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	centerWindowCmd := &cobra.Command{
		Args:  cobra.MaximumNArgs(1),
		Use:   "center [XWINDOWID]",
		Short: "Center a window",
		Long: `Center a window.
	
	If no X Window ID is given, it will center the focused window.
	By default the command does nothing to non-floating window. 
	Use the force-float argument to make the window float before centering it.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			forceFloat, err := cmd.Flags().GetBool("force-float")
			if err != nil {
				err = eris.Wrapf(internal.InvalidArgumentError, "Unable to retrieve the force-float argument")

				return
			}

			var wm *windowmgt.WindowManager
			wm, err = windowmgt.GetWindowManagerInstance()
			if err != nil {
				err = eris.Wrapf(err, "Unable to snap the windows: cannot retrieved windowmanager")

				return
			}

			internal.NormalLogger.Debug("Center command invoked.")
			var xWinIDStr *string
			if len(args) == 1 { //nolint:gomnd //Obvious number of positional arguments
				xWinIDStr = &args[0]
			}

			filterStrs, err := cmd.Flags().GetStringArray("filter")
			if err != nil {
				// TODO: deal with error
				return
			}
			windowsToProcess, err := GetWindowsToProcess(xWinIDStr, wm.GetWindows(func(w *windowmgt.Window) bool { return true }), filterStrs)

			// if len(filterStrs) == 0 || err != nil {

			// 	var targetedXWinIDPtr *xproto.Window
			// 	if len(args) == 1 { //nolint:gomnd //Obvious value to check the number of positional arguments
			// 		xWinIDArg := args[0]
			// 		var xWinIDVal uint32

			// 		if xWinIDVal, err = internal.HexStringToInt(xWinIDArg); err != nil {
			// 			err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

			// 			return
			// 		}
			// 		targetedXWinID := xproto.Window(xWinIDVal)
			// 		targetedXWinIDPtr = &targetedXWinID
			// 	}
			// 	var targetedWin *windowmgt.Window
			// 	targetedWin, err = hmi.GetTargetedWindow(targetedXWinIDPtr)
			// 	if err != nil {
			// 		if targetedXWinIDPtr == nil {
			// 			err = eris.Wrapf(err, "Unable to obtain the focused window")
			// 		} else {
			// 			err = eris.Wrapf(err, "Unable to obtain the targeted window '%#x'", *targetedXWinIDPtr)
			// 		}

			// 		return
			// 	}
			// 	windowsToCenter.Add(targetedWin)
			// } else {
			// 	var wm *windowmgt.WindowManager
			// 	wm, err = windowmgt.GetWindowManagerInstance()
			// 	if err != nil {
			// 		err = eris.Wrapf(err, "Unable to snap the windows: cannot retrieved windowmanager")

			// 		return
			// 	}
			// 	windowsToCenter, err = FilterWindowsToProcess(wm.GetWindows(func(w *windowmgt.Window) bool { return true }), filterStrs)
			// 	if err != nil {
			// 		err = eris.Wrapf(err, "Unable to snap the windows: cannot filter windows")

			// 		return
			// 	}
			// }

			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
				curWin := winIt.Value().(*windowmgt.Window)
				if err = windowmgt.CenterWindow(curWin, forceFloat); err != nil {
					err = eris.Wrapf(err, "Unable center the window '%s'", curWin.Name())
				}
			}

			return
		},
	}
	centerWindowCmd.Flags().Bool("force-float", false, "Make the window float 🤡")
	centerWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")
	parentCmd.AddCommand(centerWindowCmd)
}

func appendCenterGroupCmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	snapGroupCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
		Use:   "centergroup <GROUPID>",
		Short: "Center a group of windows",
		Long:  `TODO`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			groupID, err := strconv.Atoi(args[1])

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
				if eris.Is(err, internal.NotExistsError) {
					err = eris.Wrapf(err, "No group with ID '%d' does not exists", groupID)
				} else {
					err = eris.Wrapf(err, "Unable to retrieve the group with ID '%d'", groupID)
				}

				return
			}

			filterStrs, err := cmd.Flags().GetStringArray("filter")
			if err != nil {
				// TODO: deal with error
				return
			}
			windowsToProcess, err := FilterWindowsToProcess(group.GetWindows(), filterStrs)

			// windowsToSnap := group.GetWindows()
			// if len(filterStrs) != 0 {

			// 	windowsToSnap, err = FilterWindowsToProcess(windowsToSnap, filterStrs)
			// 	if err != nil {
			// 		err = eris.Wrapf(err, "Unable to snap the windows: cannot filter windows")

			// 		return
			// 	}
			// }

			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
				curWin := winIt.Value().(*windowmgt.Window)
				err = windowmgt.CenterWindow(curWin, forceFloat)
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
}

func AppendCenterCmd(parentCmd *cobra.Command) {
	centerCmd := &cobra.Command{
		Use:   "center",
		Short: "Center elements ",
	}
	appendCenterWindowCmd(centerCmd)
	appendCenterGroupCmd(centerCmd)
	parentCmd.AddCommand(centerCmd)
}
