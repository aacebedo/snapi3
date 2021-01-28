package hmi

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindGridDefinitionFlags(parentCmd *cobra.Command) (err error) {
	if err = viper.BindPFlag("grid.cols", parentCmd.Flags().Lookup("cols")); err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to bind the cols argument to the 'grid.cols' configuration element")

		return
	}

	if err = viper.BindPFlag("grid.rows", parentCmd.Flags().Lookup("rows")); err != nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Unable to bind the rows argument to the 'grid.rows' configuration element")

		return
	}

	return
}

func AddGridDefinitionFlags(parentCmd *cobra.Command) {
	parentCmd.Flags().Uint32P("cols", "c", 1, "The number of cols of the screen grid")
	parentCmd.Flags().Uint32P("rows", "r", 1, "The number of rows of the screen grid")
}

func GetTargetedWindow(xWinIDStr *string) (res *windowmgt.Window, err error) {
	var wm *windowmgt.WindowManager

	wm, err = windowmgt.GetWindowManagerInstance()
	if err != nil {
		err = eris.Wrapf(err, "Unable to obtain the window manager instance")

		return
	}

	if xWinIDStr != nil {
		var xWinIDVal uint32
		xWinIDVal, err = internal.HexStringToInt(*xWinIDStr)
		if err != nil {
			err = eris.Wrapf(err, "Unable to convert '%s' into an xwindow id", xWinIDStr)

			return
		}

		res, err = wm.GetWindow(xproto.Window(xWinIDVal))
		if err != nil {
			err = eris.Wrap(err, "Unable to retrieve the targeted window")

			return
		}
	} else {
		res, err = wm.GetFocusedWindow()
		if err != nil {
			err = eris.Wrap(err, "Unable to retrieve the targeted node")

			return
		}
	}

	return
}
