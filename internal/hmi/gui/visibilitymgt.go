package gui

import (
	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func appendShowGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	showGroupGUICmd := &cobra.Command{
		Use:   "show",
		Short: "Show the windows of a selected group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group to show")

				return
			}

			wins := selectedGroup.GetWindows()
			for winIt := wins.Iterator(); winIt.Next(); {
				win := winIt.Value().(*windowmgt.Window)
				if win.Show() != nil {
					internal.NormalLogger.Warnf("The window '%s' cannot been hidden", win)
				}
			}

			return
		},
	}
	parentCmd.AddCommand(showGroupGUICmd)
}

func appendHideGroupGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	showGroupGUICmd := &cobra.Command{
		Use:   "hide",
		Short: "Hide the windows of a selected group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group to remove")

				return
			}

			wins := selectedGroup.GetWindows()
			for winIt := wins.Iterator(); winIt.Next(); {
				win := winIt.Value().(*windowmgt.Window)
				if win.Hide() != nil {
					internal.NormalLogger.Warnf("The window '%s' cannot been shown", win)
				}
			}

			return
		},
	}
	parentCmd.AddCommand(showGroupGUICmd)
}

func appendToggleGroupVisibilityGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	showGroupGUICmd := &cobra.Command{
		Use:   "toggle",
		Short: "Toggle the visibility of windows of a selected group",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedGroup, err := AskUserToSelectGroup(func(group *windowmgt.WindowGroup) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the group to remove")

				return
			}

			wins := selectedGroup.GetWindows()
			for winIt := wins.Iterator(); winIt.Next(); {
				win := winIt.Value().(*windowmgt.Window)
				if win.ToggleVisibility() != nil {
					internal.NormalLogger.Warnf("Visibility of window '%s' cannot been toggled", win)
				}
			}

			return
		},
	}
	parentCmd.AddCommand(showGroupGUICmd)
}

func appendShowWindowGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	showWindowGUICmd := &cobra.Command{
		Use:   "show",
		Short: "Show the selected window",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedWin, err := AskUserToSelectWindow(func(win *windowmgt.Window) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the window to show")

				return
			}

			if selectedWin.Show() != nil {
				internal.NormalLogger.Warnf("The window '%s' cannot been shown", selectedWin)
			}

			return
		},
	}

	parentCmd.AddCommand(showWindowGUICmd)
}

func appendHideWindowGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	hideWindowGUICmd := &cobra.Command{
		Use:   "hide",
		Short: "Hide the selected window",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedWin, err := AskUserToSelectWindow(func(win *windowmgt.Window) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the window to show")

				return
			}

			if selectedWin.ToggleVisibility() != nil {
				internal.NormalLogger.Warnf("The window '%s' cannot been hidden", selectedWin)
			}

			return
		},
	}

	parentCmd.AddCommand(hideWindowGUICmd)
}

func appendToggleWindowVisibiiltyGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	showWindowGUICmd := &cobra.Command{
		Use:   "toggle",
		Short: "Toggle the visibility of the selected",
		Long:  `TODO`,

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			selectedWin, err := AskUserToSelectWindow(func(win *windowmgt.Window) bool { return true })
			if err != nil {
				err = eris.Wrap(err, "Unable to obtain the window to show")

				return
			}

			if selectedWin.Show() != nil {
				internal.NormalLogger.Warnf("Visibility of window '%s' cannot been toggled", selectedWin)
			}

			return
		},
	}

	parentCmd.AddCommand(showWindowGUICmd)
}

func AppendVisibilityMgtGUICmd(parentCmd *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	visibilityMgtGUICmd := &cobra.Command{
		Use:   "visibilitymgt",
		Short: "Hide or show windows and window groups",
	}

	windowVisibilityMgtCmd := &cobra.Command{
		Use:   "window",
		Short: "Show or hide window",
	}

	groupVisibilityMgtCmd := &cobra.Command{
		Use:   "group",
		Short: "Show or hide group",
	}
	appendHideGroupGUICmd(groupVisibilityMgtCmd)
	appendShowGroupGUICmd(groupVisibilityMgtCmd)
	appendToggleGroupVisibilityGUICmd(groupVisibilityMgtCmd)

	appendHideWindowGUICmd(windowVisibilityMgtCmd)
	appendShowWindowGUICmd(windowVisibilityMgtCmd)
	appendToggleWindowVisibiiltyGUICmd(windowVisibilityMgtCmd)

	visibilityMgtGUICmd.AddCommand(windowVisibilityMgtCmd)
	visibilityMgtGUICmd.AddCommand(groupVisibilityMgtCmd)

	parentCmd.AddCommand(visibilityMgtGUICmd)
}
