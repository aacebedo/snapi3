package cli

// import (
// 	"fmt"
// 	"strconv"
//
//
//
//
//
//
//
//
//
//
//

// 	"github.com/BurntSushi/xgb/xproto"
// 	"github.com/aacebedo/snapi3/internal"
// 	"github.com/aacebedo/snapi3/internal/hmi"
// 	"github.com/aacebedo/snapi3/internal/windowmgt"
// 	"github.com/emirpasic/gods/sets/treeset"
// 	"github.com/rotisserie/eris"
// 	"github.com/spf13/cobra"
// )

// func appendHideGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	hideGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "hide <GROUPID>",
// 		Short: "Hide a group of windows",
// 		Long:  `Hide a group of windows`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group ID '%d' is invalid", groupID)

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()

// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to find the group with ID '%d'", groupID)

// 				return
// 			}

// 			wins := group.GetWindows()
// 			for winIt := wins.Iterator(); winIt.Next(); {
// 				win := winIt.Value().(*windowmgt.Window)
// 				if hideErr := win.Hide(); hideErr != nil {
// 					internal.NormalLogger.Warnf("Unable to hide window '%s'", win)
// 					err = hideErr
// 				}
// 			}

// 			return
// 		},
// 	}
// 	parentCmd.AddCommand(hideGroupCmd)
// }

// func appendShowGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	showGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "show <GROUPID>",
// 		Short: "Show a group of windows",
// 		Long:  `Show a group of windows`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group ID '%d' is invalid", groupID)

// 				return
// 			}
// 			wgm := windowmgt.GetWindowGroupManagerInstance()

// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to find the group with ID '%d'", groupID)

// 				return
// 			}

// 			wins := group.GetWindows()
// 			for winIt := wins.Iterator(); winIt.Next(); {
// 				win := winIt.Value().(*windowmgt.Window)
// 				if showErr := win.Show(); showErr != nil {
// 					internal.NormalLogger.Warnf("Unable to show window '%s': '%s'", showErr)
// 					err = showErr
// 				}
// 			}

// 			return
// 		},
// 	}
// 	parentCmd.AddCommand(showGroupCmd)
// }

// func appendToggleGroupVisibilityCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	showGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "toggle <GROUPNAME>",
// 		Short: "Toggle visibility of a group of windows",
// 		Long:  `Toggle visibility of a group of windows`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group ID '%d' is invalid", groupID)

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()

// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to find the group with ID '%d'", groupID)

// 				return
// 			}

// 			wins := group.GetWindows()
// 			for winIt := wins.Iterator(); winIt.Next(); {
// 				win := winIt.Value().(*windowmgt.Window)
// 				if err = win.ToggleVisibility(); err != nil {
// 					err = eris.Wrapf(err, "Unable to toggle visibility of window '%s'", win)

// 					return
// 				}
// 			}

// 			return
// 		},
// 	}
// 	parentCmd.AddCommand(showGroupCmd)
// }

// func appendHideWindowCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	hideWindowCmd := &cobra.Command{
// 		Args:  cobra.RangeArgs(0, 1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "hide <WINDOWID>",
// 		Short: "Hide a window",
// 		Long:  `Hide a group`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			windowsToProcess := treeset.NewWith(windowmgt.WindowComparator)
// 			filterStrs, err := cmd.Flags().GetStringArray("filter")

// 			if len(filterStrs) == 0 || err != nil {
// 				var targetedXWinIDPtr *xproto.Window
// 				if len(args) == 1 { //nolint:gomnd //Obvious number of positional arguments
// 					xWinIDArg := args[0]
// 					var xWinIDVal uint32
// 					xWinIDVal, err = internal.HexStringToInt(xWinIDArg)
// 					if err != nil {
// 						err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

// 						return
// 					}
// 					targetedXWinID := xproto.Window(xWinIDVal)
// 					targetedXWinIDPtr = &targetedXWinID
// 				}

// 				var targetedWin *windowmgt.Window
// 				targetedWin, err = hmi.GetTargetedWindow(targetedXWinIDPtr)
// 				if err != nil {
// 					if targetedXWinIDPtr == nil {
// 						err = eris.Wrapf(err, "Unable to obtain the focused window")
// 					} else {
// 						err = eris.Wrapf(err, "Unable to obtain the targeted window '%#x'", *targetedXWinIDPtr)
// 					}

// 					return
// 				}
// 				windowsToProcess.Add(targetedWin)
// 			} else {
// 				var wm *windowmgt.WindowManager
// 				wm, err = windowmgt.GetWindowManagerInstance()
// 				if err != nil {
// 					err = eris.Wrapf(err, "Unable to snap the windows: cannot retrieved windowmanager")

// 					return
// 				}
// 				windowsToProcess, err = FilterWindowsToProcess(wm.GetWindows(func(w *windowmgt.Window) bool { return true }), filterStrs)
// 				if err != nil {
// 					err = eris.Wrapf(err, "Unable to snap the windows: cannot filter windows")

// 					return
// 				}
// 			}
// 			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
// 				curWin := winIt.Value().(*windowmgt.Window)
// 				fmt.Println(curWin.Class())
// 				if err = curWin.Hide(); err != nil {
// 					err = eris.Wrapf(err, "Unable to hide window '%s'", curWin)
// 				}
// 			}

// 			return
// 		},
// 	}

// 	hideWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")

// 	parentCmd.AddCommand(hideWindowCmd)
// }

// func appendShowWindowCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	showWindowCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "show <WINDOWID>",
// 		Short: "Show a window",
// 		Long:  `Show a window`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			var targetedXWinIDPtr *xproto.Window

// 			xWinIDArg := args[0]
// 			var xWinIDVal uint32
// 			xWinIDVal, err = internal.HexStringToInt(xWinIDArg)
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

// 				return
// 			}
// 			targetedXWinID := xproto.Window(xWinIDVal)
// 			targetedXWinIDPtr = &targetedXWinID

// 			targetedWin, err := hmi.GetTargetedWindow(targetedXWinIDPtr)
// 			if err != nil {
// 				if targetedXWinIDPtr == nil {
// 					err = eris.Wrapf(err, "Unable to obtain the focused window")
// 				} else {
// 					err = eris.Wrapf(err, "Unable to obtain the targeted window '%#x'", *targetedXWinIDPtr)
// 				}

// 				return
// 			}
// 			if err = targetedWin.Show(); err != nil {
// 				err = eris.Wrapf(err, "Unable to show window '%s'", targetedWin)
// 			}

// 			return
// 		},
// 	}

// 	showWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")

// 	parentCmd.AddCommand(showWindowCmd)
// }

// func appendToggleWindowVisibilityCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	toggleWindowCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "toggle <WINDOWID>",
// 		Short: "Toggle visibility of a window",
// 		Long:  `Toggle visibility of a window`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			var targetedXWinIDPtr *xproto.Window

// 			xWinIDArg := args[0]
// 			var xWinIDVal uint32
// 			xWinIDVal, err = internal.HexStringToInt(xWinIDArg)
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

// 				return
// 			}
// 			targetedXWinID := xproto.Window(xWinIDVal)
// 			targetedXWinIDPtr = &targetedXWinID

// 			targetedWin, err := hmi.GetTargetedWindow(targetedXWinIDPtr)
// 			if err != nil {
// 				if targetedXWinIDPtr == nil {
// 					err = eris.Wrapf(err, "Unable to obtain the focused window")
// 				} else {
// 					err = eris.Wrapf(err, "Unable to obtain the targeted window '%#x'", *targetedXWinIDPtr)
// 				}

// 				return
// 			}
// 			err = targetedWin.ToggleVisibility()
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to toggle visibilty of window '%s'", targetedWin)

// 				return
// 			}

// 			return
// 		},
// 	}
// 	toggleWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")

// 	parentCmd.AddCommand(toggleWindowCmd)
// }

// func AppendVisibilityMgtCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	visibilityMgtCmd := &cobra.Command{
// 		Use:   "visibilitymgt",
// 		Short: "Show or hide windows or groups",
// 	}

// 	windowVisibilityMgtCmd := &cobra.Command{
// 		Use:   "window",
// 		Short: "Show or hide window",
// 	}

// 	groupVisibilityMgtCmd := &cobra.Command{
// 		Use:   "group",
// 		Short: "Show or hide group",
// 	}
// 	appendHideGroupCmd(groupVisibilityMgtCmd)
// 	appendShowGroupCmd(groupVisibilityMgtCmd)
// 	appendToggleGroupVisibilityCmd(groupVisibilityMgtCmd)

// 	appendHideWindowCmd(windowVisibilityMgtCmd)
// 	appendShowWindowCmd(windowVisibilityMgtCmd)
// 	appendToggleWindowVisibilityCmd(windowVisibilityMgtCmd)

// 	visibilityMgtCmd.AddCommand(windowVisibilityMgtCmd)
// 	visibilityMgtCmd.AddCommand(groupVisibilityMgtCmd)
// 	parentCmd.AddCommand(visibilityMgtCmd)
// }
