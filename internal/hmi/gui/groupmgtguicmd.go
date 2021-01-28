package gui

import (
	"crypto/rand"

	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/hmi"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func appendAddGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	addGroupGUICmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new windows group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display groups of windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(0)

			groupName, err := AskUserToEnterNewGroup("Enter a new group", "<group name>", []string{})
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the name of the new group")

				return
			}

			wgm := windowmgt.GetWindowGroupManagerInstance()

			buf := make([]byte, 1)
			_, _ = rand.Read(buf)
			rComponent := buf[0]

			_, _ = rand.Read(buf)
			gComponent := buf[0]

			_, _ = rand.Read(buf)
			bComponent := buf[0]

			groupID, err := wgm.GenerateGroupID()
			if err != nil {
				err = eris.Wrapf(internal.InternalError, "Unable to generate an ID for new group '%s'", groupName)

				return
			}

			_, err = wgm.AddGroup(groupID, groupName, *windowmgt.NewColor(rComponent, gComponent, bComponent))
			if err != nil {
				if eris.Is(err, internal.AlreadyExistsError) {
					err = eris.Wrapf(err, "Group '%s' already exists", groupName)
				} else {
					err = eris.Wrapf(err, "Unable to add the new group '%s'", groupName)
				}

				return
			}

			return
		},
	}
	parentCmd.AddCommand(addGroupGUICmd)
}

func appendRemoveGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	removeGroupGUICmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the selected windows group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display groups of windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(viper.GetViper().GetUint32("grouplabels.displaytimeout"))

			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the name of the group to remove")

				return
			}

			wgm := windowmgt.GetWindowGroupManagerInstance()
			if err = wgm.RemoveGroup(selectedGroup.ID()); err != nil {
				if eris.Is(err, internal.NotExistsError) {
					err = eris.Wrapf(err, "The group '%s' does not exist", selectedGroup.Name())
				} else {
					err = eris.Wrapf(internal.InternalError, "Unable to remove the group '%s'", selectedGroup.Name())
				}

				return
			}

			return
		},
	}
	parentCmd.AddCommand(removeGroupGUICmd)
}

func appendAddWindowGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	addWindowGUICmd := &cobra.Command{
		Use:   "addto",
		Short: "Add a new window to a windows group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display groups of windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(viper.GetViper().GetUint32("grouplabels.displaytimeout"))

			focusedWin, err := hmi.GetTargetedWindow(nil)
			if err != nil {
				err = eris.Wrap(err, "Unable to retrieve the focused window")

				return
			}

			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool {
				return !group.Contains(focusedWin)
			})
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group to which the window shall be added")

				return
			}

			if err = selectedGroup.AddWindow(focusedWin); err != nil {
				if eris.Is(err, internal.AlreadyExistsError) {
					err = eris.Wrapf(err, "Window '%s' already belongs to group '%s'", focusedWin.Name(), selectedGroup.Name())
				} else {
					err = eris.Wrapf(internal.InternalError, "Unable to add a the window '%s' to group '%s'", focusedWin.Name(), selectedGroup.Name())
				}

				return
			}

			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display group labels on windows")

				return
			}

			return
		},
	}
	parentCmd.AddCommand(addWindowGUICmd)
}

func appendRemoveWindowGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	removeWindowGUICmd := &cobra.Command{
		Use:   "removefrom",
		Short: "Remove a window from a windows group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display group labels on windows")

				return
			}
			defer windowmgt.HideAllWindowLabels(viper.GetViper().GetUint32("grouplabels.displaytimeout"))

			focusedWin, err := hmi.GetTargetedWindow(nil)
			if err != nil {
				err = eris.Wrap(err, "Unable to retrieve the focused window")

				return
			}

			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool {
				return group.Contains(focusedWin)
			})
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group from which the window shall be removed")

				return
			}

			if err = selectedGroup.RemoveWindow(focusedWin); err != nil {
				if eris.Is(err, internal.NotExistsError) {
					err = eris.Wrapf(err, "Window '%s' does not belong to group '%s'", focusedWin.Name(), selectedGroup.Name())
				} else {
					err = eris.Wrapf(err, "Unable to remove window '%s' from group '%s'", focusedWin.Name(), selectedGroup.Name())
				}

				return
			}

			if err = windowmgt.ShowAllWindowLabels(); err != nil {
				err = eris.Wrap(err, "Unable to display group labels on windows")

				return
			}

			return
		},
	}

	parentCmd.AddCommand(removeWindowGUICmd)
}

func AppendGroupMgtGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	groupMgtGUICmd := &cobra.Command{
		Use:   "groupmgt",
		Short: "Manage the windows group",
	}

	groupMgtGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Add or remove groups",
	}

	groupMgtWindowCmd := &cobra.Command{
		Use:   "window",
		Short: "Add or remove windows to groups",
	}

	appendAddGroupGUICmd(groupMgtGroupCmd)
	appendRemoveGroupGUICmd(groupMgtGroupCmd)

	appendAddWindowGUICmd(groupMgtWindowCmd)
	appendRemoveWindowGUICmd(groupMgtWindowCmd)

	groupMgtGUICmd.AddCommand(groupMgtGroupCmd)
	groupMgtGUICmd.AddCommand(groupMgtWindowCmd)

	parentCmd.AddCommand(groupMgtGUICmd)
}
