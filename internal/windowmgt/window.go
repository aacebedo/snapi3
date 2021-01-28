package windowmgt

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/aacebedo/snapi3/internal"
	"github.com/emirpasic/gods/utils"
	"github.com/rotisserie/eris"
	"go.i3wm.org/i3/v4"
)

func WindowComparator(left, right interface{}) int {
	leftWin := left.(*Window)
	rightWin := right.(*Window)

	return utils.UInt32Comparator(uint32(leftWin.XWinID()), uint32(rightWin.XWinID()))
}

type Geometry struct {
	X           int32
	Y           int32
	Width       uint32
	Height      uint32
	BorderWidth uint32
}

func NewGeometry(x, y int32, width, height, borderWidth uint32) (res *Geometry) {
	return &Geometry{X: x, Y: y, Width: width, Height: height, BorderWidth: borderWidth}
}

type Window struct {
	xWinID xproto.Window
	xUtil  *xgbutil.XUtil
}

func NewWindow(xWinID xproto.Window, xUtil *xgbutil.XUtil) (res *Window) {
	res = &Window{xWinID: xWinID, xUtil: xUtil}

	return
}

func (win *Window) XWinID() (res xproto.Window) {
	return win.xWinID
}

func (win *Window) State() (res *icccm.WmState) {
	res, _ = icccm.WmStateGet(win.xUtil, win.xWinID)

	return
}
func (win *Window) Name() (res string) {
	res, _ = icccm.WmNameGet(win.xUtil, win.xWinID)

	return
}

func (win *Window) Class() (res string) {
	class, _ := icccm.WmClassGet(win.xUtil, win.xWinID)

	res = class.Class

	return
}
func (win *Window) Types() (res []string) {
	res, _ = ewmh.WmWindowTypeGet(win.xUtil, win.xWinID)

	return
}

func (win *Window) String() (res string) {
	return fmt.Sprintf("%s (%#x)", win.Name(), win.xWinID)
}

func (win *Window) Resize(width, height uint32) (err error) {
	if err = ewmh.ResizeWindow(win.xUtil, win.xWinID, int(width), int(height)); err != nil {
		err = eris.Wrap(err, "Unable to resize the window")

		return
	}

	internal.NormalLogger.Debugf("Resized the window '%d' to width='%d', height='%d'.", win.xWinID, height, width)

	return
}

func (win *Window) Move(x, y int32) (err error) {
	if err = ewmh.MoveWindow(win.xUtil, win.xWinID, int(x), int(y)); err != nil {
		err = eris.Wrap(err, "Unable to move the window")

		return
	}

	internal.NormalLogger.Debugf("Moved the window '%d' to x='%d', y='%d')", win.xWinID, x, y)

	return
}

func (win *Window) MoveResize(x, y int32, width, height uint32) (err error) {
	if err = ewmh.MoveresizeWindow(win.xUtil, win.xWinID, int(x), int(y), int(width), int(height)); err != nil {
		err = eris.Errorf("Unable to move and resize the window: %s", err)

		return
	}

	internal.NormalLogger.Debugf("Moved and resized the window '%d' to x='%d', y='%d' and resized to width='%d', height='%d').",
		win.xWinID, x, y, height, width)

	return
}

func (win *Window) IsVisible() (res bool, err error) {
	res = false
	windowAttrReq := xproto.GetWindowAttributes(win.xUtil.Conn(), win.xWinID)

	windowAttr, err := windowAttrReq.Reply()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to obtain attributes of window '%s'", win)

		return
	}

	res = (windowAttr.MapState == xproto.MapStateViewable)

	return
}

func (win *Window) SetAtom(atomType xproto.Atom, atomName string, atomFormat byte, atomValue []byte, atomValueLen uint32) (err error) {
	internAtom, err := xproto.InternAtom(win.xUtil.Conn(), false, uint16(len(atomName)), atomName).Reply()
	if err != nil {
		internal.VerboseLogger.Errorf("Error while setting X11 atom: %s", err)
		err = eris.Wrapf(internal.InternalError, "Unable to get the internal atom for '%s'", atomName)

		return
	}

	err = xproto.ChangePropertyChecked(win.xUtil.Conn(), xproto.PropModeReplace, win.XWinID(), internAtom.Atom, atomType, atomFormat, atomValueLen, atomValue).Check()
	if err != nil {
		internal.VerboseLogger.Errorf("Error while setting a property on the window: %s", err)
		err = eris.Wrapf(internal.InternalError, "Unable to get the internal atom for '%s'", atomName)

		return
	}

	return
}

func (win *Window) GetAtom(atomName string) (value []byte, valueLen uint32, err error) {
	internAtom, err := xproto.InternAtom(win.xUtil.Conn(), false, uint16(len(atomName)), atomName).Reply()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to get the internal atom for '%s'", atomName)

		return
	}

	prop, err := xproto.GetProperty(win.xUtil.Conn(), false, win.xWinID, internAtom.Atom, xproto.AtomAny, 0, (1<<32)-1).Reply()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to get the internal atom for '%s'", atomName)

		return
	}
	value = prop.Value
	valueLen = prop.ValueLen
	return
}

func (win *Window) Show() (err error) {
	internal.NormalLogger.Debugf("Show window '%#x''%s'", win.xWinID, win.Name())
	xwin := xwindow.New(win.xUtil, win.xWinID)
	xwin.Map()

	return
}

func (win *Window) ToggleVisibility() (err error) {
	winVisibility, err := win.IsVisible()
	if err != nil {
		err = eris.Wrapf(err, "Unable to toggle visibilty of window '%s': Failed to get visibility state", win)

		return
	}

	if winVisibility {
		if err = win.Hide(); err != nil {
			err = eris.Wrapf(err, "Unable to show window '%s'", win)
		}
	} else {
		if err = win.Show(); err != nil {
			err = eris.Wrapf(err, "Unable to show window '%s'", win)
		}
	}
	return err
}

func (win *Window) Hide() (err error) {
	internal.NormalLogger.Debugf("Hide window '%#x' '%s'", win.xWinID, win.Name())
	xwin := xwindow.New(win.xUtil, win.xWinID)
	xwin.Unmap()

	return
}

func (win *Window) GetNodeID() (res i3.NodeID, err error) {
	i3Tree, err := i3.GetTree()
	if err != nil {
		err = eris.Wrap(internal.InternalError, "Unable to get i3 tree")

		return
	}

	node := i3Tree.Root.FindChild(func(n *i3.Node) bool {
		return xproto.Window(n.Window) == win.xWinID
	})

	if node == nil {
		err = eris.Wrapf(internal.NotExistsError, "Window '%#x' does not exist in the i3 tree", win.xWinID)

		return
	}

	res = node.ID

	return
}

func (win *Window) GetGeometry() (res *Geometry, err error) {
	nodeID, err := win.GetNodeID()
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the corresponding i3 node")

		return
	}

	var wm *WindowManager

	wm, err = GetWindowManagerInstance()
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the window manager")

		return
	}

	node, err := wm.GetNode(nodeID)
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the corresponding i3 node")

		return
	}

	res = NewGeometry(int32(node.WindowRect.X), int32(node.WindowRect.Y),
		uint32(node.WindowRect.Width), uint32(node.WindowRect.Height),
		uint32(node.CurrentBorderWidth))

	return
}

func (win *Window) IsWindowFloating() (res bool, err error) {
	res = false

	nodeID, err := win.GetNodeID()
	if err != nil {
		err = eris.Wrapf(err, "Unable to obtain the floating state of the window '%s'", win)

		return
	}

	var wm *WindowManager

	wm, err = GetWindowManagerInstance()
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the window manager")

		return
	}

	floatingNodes, err := wm.GetNodes(func(n *i3.Node) bool {
		return n.Type == i3.FloatingCon && len(n.Nodes) == 1
	})
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the floating nodes")

		return
	}

	for it := floatingNodes.Iterator(); it.Next(); {
		if it.Value().(*i3.Node).Nodes[0].ID == nodeID {
			res = true

			break
		}
	}

	return res, err
}
