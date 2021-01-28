package windowmgt

import (
	"fmt"

	"github.com/aacebedo/snapi3/internal"
	"github.com/rotisserie/eris"
	"go.i3wm.org/i3/v4"
)

// func (wm *WindowManager) IsWindowFloating(win *Window) (res bool) {
// 	ancestors, err := wm.getI3NodeAncestors(win)
// 	res = false

// 	if err != nil {
// 		internal.NormalLogger.Warnf("Unable to retrieve the i3 ancestors of the window. Floating state will be invalid")

// 		return
// 	}

// 	it := ancestors.Iterator()
// 	for it.End(); it.Prev(); {
// 		if it.Value().(*i3.Node).Type == i3.FloatingCon {
// 			res = true

// 			return
// 		}
// 	}

// 	return
// }

func ensureWindowIsFloating(win *Window, forceFloat bool) (err error) {
	winIsFloating, err := win.IsWindowFloating()
	if err != nil {
		err = eris.Wrapf(err, "Unable to know if window '%s' is floating", win)

		return
	}

	if !winIsFloating {
		internal.VerboseLogger.Debugf("Window '%s' is not floating", win.Name())

		if !forceFloat {
			err = eris.Wrapf(internal.InternalError,
				"The targeted window '%s' is not floating and "+
					"the force-float flag has not been used.", win.Name())

			return
		}
	}

	if err = SetWindowFloatingState(win, true); err != nil {
		err = eris.Wrapf(err, "Unable to make the window '%s' to float", win.Name())

		return
	}

	internal.VerboseLogger.Debugf("Window '%s' is now floating", win.Name())

	return
}

func CenterWindow(win *Window, forceFloat bool) (err error) {
	if err = ensureWindowIsFloating(win, forceFloat); err != nil {
		err = eris.Wrapf(err, "Unable to ensure the window '%s' is floating", win.Name())

		return
	}

	winGeom, err := win.GetGeometry()
	if err != nil {
		err = eris.Wrapf(err, "Unable to retrieve geometry of window '%s'", win.Name())

		return
	}

	screenWidth, screenHeight, err := GetScreenSize()
	if err != nil {
		err = eris.Wrapf(err, "Unable to get the size of the screen")

		return
	}

	//nolint:gomnd //Obvious calculation to obtain the center of the screen
	newX := int32(screenWidth/2) - int32(winGeom.Width/2)
	//nolint:gomnd //Obvious calculation to obtain the center of the screen
	newY := int32(screenHeight/2) - int32(winGeom.Height/2)

	if err = win.Move(newX, newY); err != nil {
		err = eris.Wrapf(err, "Unable to move the window '%s'", win.Name())

		return
	}

	return err
}

func SnapWindow(win *Window, rows, cols, selectedRow, selectedCol uint32, forceFloat bool) (err error) {
	internal.VerboseLogger.Debugf("Snapping window '%s' with rows='%d', cols='%d', selectedRow='%d', selectedCol='%d'.",
		win.Name(), rows, cols, selectedRow, selectedCol)

	if err = ensureWindowIsFloating(win, forceFloat); err != nil {
		err = eris.Wrapf(err, "Unable to ensure the window '%s' is floating", win.Name())

		return
	}

	winGeom, err := win.GetGeometry()
	if err != nil {
		err = eris.Wrapf(err, "Unable to get geometry of window '%s'", win.Name())

		return
	}

	screenWidth, screenHeight, err := GetScreenSize()
	if err != nil {
		err = eris.Wrapf(err, "Unable to get the size of the screen")

		return
	}

	sg, err := NewScreenGrid(rows, cols, screenWidth, screenHeight)
	if err != nil {
		err = eris.Wrap(err, "Unable to create the screen grid")

		return
	}

	pos, err := sg.GetPosition(selectedRow, selectedCol)
	if err != nil {
		err = eris.Wrap(err, "Unable to obtain the requested position of the window on the grid")

		return
	}

	borderWidth := winGeom.BorderWidth
	newX := pos.X() + int32(borderWidth)
	newY := pos.Y() + int32(borderWidth)

	newWidth := sg.CellWidth() - borderWidth*2
	newHeight := sg.CellHeight() - borderWidth*2

	internal.VerboseLogger.Debugf("New position of the window '%s' will be x='%d', y='%d'.", win.Name(), newX, newY)
	internal.VerboseLogger.Debugf("New dimension of the window '%s' will be width='%d', height='%d'.", win.Name(), newWidth, newHeight)

	if err = win.MoveResize(newX, newY, newWidth, newHeight); err != nil {
		err = eris.Wrap(err, "Unable to resize the window")

		return
	}

	return err
}

func SetWindowFloatingState(win *Window, isFloating bool) (err error) {
	windowIsFloating, err := win.IsWindowFloating()
	if err != nil {
		err = eris.Wrapf(err, "Cannot make the window '%s' float", win.Name())

		return
	}

	if windowIsFloating != isFloating {
		var i3Cmd string
		if isFloating {
			i3Cmd = fmt.Sprintf("[id=%d] floating enable", win.XWinID())
		} else {
			i3Cmd = fmt.Sprintf("[id=%d] floating disable", win.XWinID())
		}

		internal.NormalLogger.Debugf("Run i3 command: '%s'", i3Cmd)

		if _, err = i3.RunCommand(i3Cmd); err != nil {
			err = eris.Wrapf(err, "Cannot make the window '%s' float", win.Name())

			return
		}
	}

	return
}

func ShowAllWindowLabels() (err error) {
	ld, err := GetLabelDrawerInstance()
	if err != nil {
		err = eris.Wrap(err, "Unable to obtain an instance of the LabelDrawer")

		return
	}

	wm, err := GetWindowManagerInstance()
	//TODO: deal with error
	wins := wm.GetWindows(func(win *Window) bool { return true })
	for winIt := wins.Iterator(); winIt.Next(); {

		curWin := winIt.Value().(*Window)

		labelErr := ld.ShowWindowLabels(curWin)
		if labelErr != nil {
			err = eris.Wrapf(err, "Unable to display window labels of window '%s'", curWin.Name())

			return
		}
	}

	return
}

func HideAllWindowLabels(displayTimeout uint32) {
	// time.Sleep(time.Duration(displayTimeout) * time.Second)
	// gtk.MainQuit()
}

func GetScreenSize() (width, height uint32, err error) {
	i3Tree, err := i3.GetTree()
	if err != nil {
		err = eris.Wrap(internal.InternalError, "Unable to retrieve the i3 tree")

		return
	}

	focusedOutput := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.OutputNode
	})
	if focusedOutput == nil {
		err = eris.Wrapf(internal.InternalError, "Unable to retrieve focused output")

		return
	}

	width = uint32(focusedOutput.Rect.Width)
	height = uint32(focusedOutput.Rect.Height)

	return
}
