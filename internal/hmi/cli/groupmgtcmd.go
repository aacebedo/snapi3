package cli

// import (
// 	"regexp"
// 	"strconv"
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
// 	"github.com/rotisserie/eris"
// 	"github.com/spf13/cobra"
// )

// func appendAddGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	addGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(2), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "add <GROUPNAME> <COLOR>",
// 		Short: "Add a new windows group",
// 		Long: `Add a new windows group that can be used to perform different tasks

// 		The name of the group and its color is needed when adding it.`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupName := args[0]
// 			matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			if !matching || err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group name '%s' is invalid", groupName)

// 				return
// 			}
// 			groupColorArg := args[1]

// 			groupColor, err := windowmgt.NewColorFromHTMLCode(groupColorArg)
// 			if err != nil {
// 				err = eris.Wrapf(err, "The HTML code '%s' cannot be parsed into a color", groupColorArg)

// 				return
// 			}

// 			if err = windowmgt.ShowAllWindowLabels(); err != nil {
// 				err = eris.Wrap(err, "Unable to display groups of windows")

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			groupID, err := wgm.GenerateGroupID()
// 			if err != nil {
// 				err = eris.Wrapf(internal.InternalError, "Unable to generate an ID for new group '%s'", groupName)
// 				return
// 			}

// 			if _, err = wgm.AddGroup(groupID, groupName, *groupColor); err != nil {
// 				if eris.Is(err, internal.AlreadyExistsError) {
// 					err = eris.Wrapf(err, "Group '%s' already exists", groupName)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to add the group '%s'", groupName)
// 				}

// 				return
// 			}

// 			return
// 		},
// 	}
// 	parentCmd.AddCommand(addGroupCmd)
// }

// func appendSetGroupColorCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	setGroupColorCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(2), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "setcolor <GROUPID> <COLOR>",
// 		Short: "Set the color of a windows group",
// 		Long:  `Set the color of a windows group`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			if err = windowmgt.ShowAllWindowLabels(); err != nil {
// 				err = eris.Wrap(err, "Unable to display groups of windows")

// 				return
// 			}
// 			defer windowmgt.HideAllWindowLabels(0)
// 			groupID, err := strconv.Atoi(args[0])
// 			/*matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			if !matching || err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group name '%s' is invalid", groupName)

// 				return
// 			}*/
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group with ID '%d' is invalid", groupID)

// 				return
// 			}
// 			groupColorArg := args[1]

// 			groupColor, err := windowmgt.NewColorFromHTMLCode(groupColorArg)
// 			if err != nil {
// 				err = eris.Wrapf(err, "The HTML code '%s' cannot be parsed into a color", groupColorArg)

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			winGroup, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "The group with ID '%d' does not exist", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to add the group with ID '%d'", groupID)
// 				}

// 				return
// 			}

// 			winGroup.SetColor(*groupColor)

// 			return
// 		},
// 	}

// 	parentCmd.AddCommand(setGroupColorCmd)
// }

// func appendRenameGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	renameGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(2), //nolint:gomnd //Obvious number of positional arguments
// 		Use:   "rename <GROUPID> <NEWGROUPNAME>",
// 		Short: "Rename a window group",
// 		Long:  `Rename a window group`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			// matching, err := regexp.MatchString("^(\\w+\\s?)$", oldGroupName)
// 			// if !matching || err != nil {
// 			// 	err = eris.Wrapf(internal.InvalidArgumentError, "Old group name '%s' is invalid", oldGroupName)

// 			// 	return
// 			// }
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group ID '%d'", groupID)

// 				return
// 			}

// 			newGroupName := args[1]
// 			matching, err := regexp.MatchString("^(\\w+\\s?)$", newGroupName)
// 			if !matching || err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "New group name '%s' is invalid", newGroupName)

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			// if wgm.Contains(uint8(groupID)) {
// 			// 	err = eris.Wrapf(internal.AlreadyExistsError, "Group '%s' already exists", newGroupName)

// 			// 	return
// 			// }

// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "Group with ID '%d' does not exists", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to obtain the group with ID '%d'", groupID)
// 				}

// 				return
// 			}

// 			group.SetName(newGroupName)
// 			// if err = wgm.RemoveGroup(oldGroup.Name()); err != nil {
// 			// 	err = eris.Wrapf(err, "Unable to remove the old group '%s'", oldGroupName)

// 			// 	return
// 			// }

// 			// newGroup, err := wgm.AddGroup(newGroupName, oldGroup.Color())
// 			// if err != nil {
// 			// 	err = eris.Wrapf(err, "Unable to add the new group '%s'", newGroupName)

// 			// 	return
// 			// }

// 			// for winIt := oldGroup.GetWindows().Iterator(); winIt.Next(); {
// 			// 	curWin := winIt.Value().(*windowmgt.Window)

// 			// 	if addErr := newGroup.AddWindow(curWin); addErr != nil {
// 			// 		internal.NormalLogger.Warnf("Unable to add the window '%s' to new group '%s'", curWin.Name(), oldGroupName)
// 			// 	}
// 			// }

// 			return
// 		},
// 	}

// 	parentCmd.AddCommand(renameGroupCmd)
// }

// func appendRemoveGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	removeGroupCmd := &cobra.Command{
// 		Args:  cobra.ExactArgs(1),
// 		Use:   "remove <GROUPID>",
// 		Short: "Remove a new windows group",
// 		Long: `Remove a new windows group that can be used to perform different tasks

// 		The name of the group and its color is needed when adding it.`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			// matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			// if !matching || err != nil {
// 			// 	err = eris.Wrapf(internal.InvalidArgumentError, "Group name '%s' is invalid", groupName)

// 			// 	return
// 			// }

// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group ID '%d' is invalid", groupID)

// 				return
// 			}
// 			wgm := windowmgt.GetWindowGroupManagerInstance()

// 			if err = wgm.RemoveGroup(windowmgt.WindowGroupIdentifier(groupID)); err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "The group with ID '%d' does not exist", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to remove the group with ID '%d'", groupID)
// 				}

// 				return
// 			}

// 			return
// 		},
// 	}
// 	parentCmd.AddCommand(removeGroupCmd)
// }

// //nolint:dupl //Function is very similar to appendRemoveWindowToGroupCmd but I prefer to keep it that way
// func appendAddWindowToGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	addWindowCmd := &cobra.Command{
// 		Args:  cobra.RangeArgs(1, 2),
// 		Use:   "add <GROUPID> [XWINDOWID]",
// 		Short: "Add a window to a windows group",
// 		Long:  `Add a window to a windows group.`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			// matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			// if !matching || err != nil {
// 			// 	err = eris.Wrapf(internal.InvalidArgumentError, "Group name is '%s'", groupName)

// 			// 	return
// 			// }
// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Group ID '%d' is invalid", groupID)

// 				return
// 			}

// 			var wm *windowmgt.WindowManager
// 			wm, err = windowmgt.GetWindowManagerInstance()
// 			if err != nil {
// 				err = eris.Wrapf(err, "Unable to snap the windows: cannot retrieved windowmanager")

// 				return
// 			}

// 			var xWinIDStr *string
// 			if len(args) == 2 { //nolint:gomnd //Obvious number of positional arguments
// 				xWinIDStr = &args[1]
// 			}

// 			filterStrs, err := cmd.Flags().GetStringArray("filter")
// 			if err != nil {
// 				// TODO: deal with error
// 				return
// 			}

// 			windowsToProcess, err := GetWindowsToProcess(xWinIDStr, wm.GetWindows(func(w *windowmgt.Window) bool { return true }), filterStrs)

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "Group with ID '%d' does not exist", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to find the group with ID '%d'", groupID)
// 				}

// 				return
// 			}

// 			for winIt := windowsToProcess.Iterator(); winIt.Next(); {
// 				curWin := winIt.Value().(*windowmgt.Window)
// 				if err = group.AddWindow(curWin); err != nil {
// 					if eris.Is(err, internal.AlreadyExistsError) {
// 						err = eris.Wrapf(err, "Window '%s' already belongs to group '%s'", curWin.Name(), group.Name())
// 					} else {
// 						err = eris.Wrapf(err, "Unable to add window '%s' to group '%s'", curWin.Name(), group.Name())
// 					}

// 					return
// 				}
// 			}

// 			return
// 		},
// 	}

// 	addWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")
// 	parentCmd.AddCommand(addWindowCmd)
// }

// //nolint:dupl //Function is very similar to appendAddWindowToGroupCmd but I prefer to keep it that way
// func appendRemoveWindowToGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	removeWindowCmd := &cobra.Command{
// 		Args:  cobra.RangeArgs(1, 2),
// 		Use:   "remove <GROUPID> [XWINDOWID]",
// 		Short: "Remove a window from a windows group",
// 		Long:  `Remove a window from a windows group.`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			// matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			// if !matching || err != nil {
// 			// 	err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group name '%s'", groupName)

// 			// 	return
// 			// }

// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group ID '%d'", groupID)

// 				return
// 			}

// 			var targetedXWinIDPtr *xproto.Window
// 			if len(args) == 2 { //nolint:gomnd //Obvious number of positional arguments
// 				xWinIDArg := args[1]
// 				var xWinIDVal uint32
// 				xWinIDVal, err = internal.HexStringToInt(xWinIDArg)
// 				if err != nil {
// 					err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

// 					return
// 				}
// 				targetedXWinID := xproto.Window(xWinIDVal)
// 				targetedXWinIDPtr = &targetedXWinID
// 			}

// 			targetedWin, err := hmi.GetTargetedWindow(targetedXWinIDPtr)
// 			if err != nil {
// 				if targetedXWinIDPtr == nil {
// 					err = eris.Wrapf(err, "Unable to obtain the focused window")
// 				} else {
// 					err = eris.Wrapf(err, "Unable to obtain the targeted window '%#x'", *targetedXWinIDPtr)
// 				}

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "Group with ID '%d' does not exist", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to find the group '%s'", groupID)
// 				}

// 				return
// 			}

// 			if err = group.RemoveWindow(targetedWin); err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "Window '%s' does not belong to group '%s'", targetedWin.Name(), group.Name())
// 				} else {
// 					err = eris.Wrapf(err, "Unable to remove window '%s' from group '%s'", targetedWin.Name(), group.Name())
// 				}

// 				return
// 			}

// 			return
// 		},
// 	}

// 	removeWindowCmd.Flags().StringArrayP("filter", "f", []string{}, "Apply only to windows matching this filter")
// 	parentCmd.AddCommand(removeWindowCmd)
// }

// //nolint:dupl //Function is very similar to appendAddWindowToGroupCmd but I prefer to keep it that way
// func appendToggleWindowToGroupCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	toggleWindowCmd := &cobra.Command{
// 		Args:  cobra.RangeArgs(1, 2),
// 		Use:   "togglewindow <GROUPID> [XWINDOWID]",
// 		Short: "Toggle a window in or out a windows group",
// 		Long:  `Toggle a window in or out a windows group.`,
// 		RunE: func(cmd *cobra.Command, args []string) (err error) {
// 			groupID, err := strconv.Atoi(args[0])
// 			// matching, err := regexp.MatchString("^(\\w+\\s?)$", groupName)
// 			// if !matching || err != nil {
// 			// 	err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group name '%s'", groupName)

// 			// 	return
// 			// }

// 			if err != nil {
// 				err = eris.Wrapf(internal.InvalidArgumentError, "Invalid group ID '%d'", groupID)

// 				return
// 			}

// 			var targetedXWinIDPtr *xproto.Window
// 			if len(args) == 2 { //nolint:gomnd //Obvious number of positional arguments
// 				xWinIDArg := args[1]
// 				var xWinIDVal uint32
// 				xWinIDVal, err = internal.HexStringToInt(xWinIDArg)
// 				if err != nil {
// 					err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDArg)

// 					return
// 				}
// 				targetedXWinID := xproto.Window(xWinIDVal)
// 				targetedXWinIDPtr = &targetedXWinID
// 			}

// 			targetedWin, err := hmi.GetTargetedWindow(targetedXWinIDPtr)
// 			if err != nil {
// 				if targetedXWinIDPtr == nil {
// 					err = eris.Wrapf(err, "Unable to obtain the focused window")
// 				} else {
// 					err = eris.Wrapf(err, "Unable to obtain the targeted window '0x%#X'", *targetedXWinIDPtr)
// 				}

// 				return
// 			}

// 			wgm := windowmgt.GetWindowGroupManagerInstance()
// 			group, err := wgm.GetGroup(windowmgt.WindowGroupIdentifier(groupID))
// 			if err != nil {
// 				if eris.Is(err, internal.NotExistsError) {
// 					err = eris.Wrapf(err, "Group with ID '%d' does not exist", groupID)
// 				} else {
// 					err = eris.Wrapf(err, "Unable to find the group with ID '%d'", groupID)
// 				}

// 				return
// 			}

// 			if !group.Contains(targetedWin) {
// 				if err = group.RemoveWindow(targetedWin); err != nil {
// 					err = eris.Wrapf(err, "Unable to remove window '0x%#X' from group '%s'", targetedWin.XWinID(), group.Name())

// 					return
// 				}
// 			} else {
// 				if err = group.AddWindow(targetedWin); err != nil {
// 					err = eris.Wrapf(err, "Unable to add window '0x%#X' from group '%s'", targetedWin.XWinID(), group.Name())

// 					return
// 				}
// 			}

// 			return
// 		},
// 	}

// 	parentCmd.AddCommand(toggleWindowCmd)
// }

// func AppendGroupMgtCmd(parentCmd *cobra.Command) {
// 	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
// 	groupMgtCmd := &cobra.Command{
// 		Use:   "groupmgt",
// 		Short: "Manage groups of windows",
// 	}

// 	appendAddGroupCmd(groupMgtCmd)
// 	appendRemoveGroupCmd(groupMgtCmd)
// 	appendSetGroupColorCmd(groupMgtCmd)
// 	appendRenameGroupCmd(groupMgtCmd)

// 	appendAddWindowToGroupCmd(groupMgtCmd)
// 	appendRemoveWindowToGroupCmd(groupMgtCmd)
// 	appendToggleWindowToGroupCmd(groupMgtCmd)

// 	parentCmd.AddCommand(groupMgtCmd)
// }
