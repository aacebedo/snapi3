package windowmgt

import (
	"sync"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/aacebedo/snapi3/internal"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"github.com/rotisserie/eris"
	"go.i3wm.org/i3/v4"
)

func NodeComparator(left, right interface{}) int {
	leftNode := left.(*i3.Node)
	rightNode := right.(*i3.Node)

	switch {
	case leftNode.ID > rightNode.ID:
		return 1
	case leftNode.ID < rightNode.ID:
		return -1
	default:
		return 0
	}
}

type WindowManager struct {
	windows  treeset.Set
	winMutex sync.Mutex
	xUtil    *xgbutil.XUtil
}

var (
	wmInstance     *WindowManager  //nolint:gochecknoglobals //Needed for singleton pattern
	wmInstanceLock = &sync.Mutex{} //nolint:gochecknoglobals //Needed for singleton pattern
)

func WindowIDComparator(left, right interface{}) int {
	leftWin := left.(xproto.Window)
	rightWin := right.(xproto.Window)

	return utils.UInt32Comparator(uint32(leftWin), uint32(rightWin))
}

func GetWindowManagerInstance() (res *WindowManager, err error) {
	if wmInstance == nil {
		wmInstanceLock.Lock()
		defer wmInstanceLock.Unlock()

		if wmInstance == nil {
			var xUtil *xgbutil.XUtil

			xUtil, err = xgbutil.NewConn()
			if err != nil {
				err = eris.Wrapf(internal.InternalError, "Unable to create a xgb connection")

				return
			}

			wmInstance = &WindowManager{windows: *treeset.NewWith(WindowComparator), xUtil: xUtil}
		}
	}

	return wmInstance, err
}

func (wm *WindowManager) getI3NodeAncestors(win *Window) (res *arraylist.List, err error) {
	res = arraylist.New()

	i3Tree, err := i3.GetTree()
	if err != nil {
		err = eris.Wrap(internal.InternalError, "Unable to get i3 tree")

		return
	}

	node := i3Tree.Root.FindChild(func(n *i3.Node) bool {
		return xproto.Window(n.Window) == win.xWinID
	})
	if node == nil {
		err = eris.Wrapf(internal.AlreadyExistsError, "Unable to find an i3 node with the a xwindow equal to '%d'", win.xWinID)

		return
	}

	curNodeID, err := win.GetNodeID()
	if err != nil {
		err = eris.Wrapf(err, "Unable to obtain i3 nodeUID attached to window '%s'", win)

		return
	}

	for {
		parentNode := i3Tree.Root.FindChild(func(n *i3.Node) bool {
			childrens := arraylist.New()
			for _, n := range n.Nodes {
				childrens.Add(n.ID)
			}
			for _, n := range n.FloatingNodes {
				childrens.Add(n.ID)
			}

			return childrens.Contains(curNodeID)
		})
		if parentNode != nil {
			res.Insert(0, parentNode)
			curNodeID = parentNode.ID
		} else {
			return
		}
	}
}

func (wm *WindowManager) AddWindow(xWinID xproto.Window) (res *Window, err error) {
	wm.winMutex.Lock()
	defer wm.winMutex.Unlock()
	idx, winItf := wm.windows.Find(func(index int, value interface{}) bool {
		curWin := value.(*Window)

		return curWin.XWinID() == xWinID
	})

	if idx == -1 {
		res = NewWindow(xWinID, wm.xUtil)

		wm.windows.Add(res)
		internal.VerboseLogger.Debugf("Window '%s' is managed", res)
	} else {
		win := winItf.(*Window)
		err = eris.Wrapf(internal.AlreadyExistsError, "Window with node id '%s' already managed", win)

		return
	}

	return
}

func (wm *WindowManager) RemoveWindow(xWinID xproto.Window) (err error) {
	wm.winMutex.Lock()
	defer wm.winMutex.Unlock()
	idx, winToRemove := wm.windows.Find(func(index int, value interface{}) bool {
		curWin := value.(*Window)

		return curWin.XWinID() == xWinID
	})

	if idx != -1 {
		wm.windows.Remove(winToRemove)
		internal.NormalLogger.Debugf("Window '%s' is unmanaged", winToRemove.(*Window))
	} else {
		err = eris.Wrapf(internal.NotExistsError, "Window '%#x' not managed", winToRemove.(*Window).XWinID())

		return
	}

	return
}

func (wm *WindowManager) GetWindows(predicate func(*Window) bool) (res *treeset.Set) {
	wm.winMutex.Lock()
	defer wm.winMutex.Unlock()
	res = wm.windows.Select(func(index int, value interface{}) bool {
		curWin := value.(*Window)

		return predicate(curWin)
	})

	return
}

func (wm *WindowManager) GetWindow(xWinID xproto.Window) (res *Window, err error) {
	winSet := wm.GetWindows(func(win *Window) bool {
		return win.XWinID() == xWinID
	})

	if winSet.Size() != 1 {
		err = eris.Wrapf(internal.NotExistsError, "Unable to find the window with ID '%d'", xWinID)

		return
	}

	res = winSet.Values()[0].(*Window)

	return
}

func (wm *WindowManager) isWindowManageable(xWinID xproto.Window) (res bool, err error) {
	res = false
	win := NewWindow(xWinID, wm.xUtil)

	winState := win.State()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to retrieve the state of the window '%s'", win)

		return
	}
	//internal.VerboseLogger.Warnf("State of window '%#x' is '%d'", xWinID, winState)

	windowTypes := win.Types()
	//internal.VerboseLogger.Debug("Types of window '%#x' is '%s'", xWinID, windowTypes)

	winTypeIsManageable := false

	for _, v := range windowTypes {
		if v == "_NET_WM_WINDOW_TYPE_NORMAL" || v == "_NET_WM_WINDOW_TYPE_UTILITY" {
			winTypeIsManageable = true

			break
		}
	}

	windowAttrReq := xproto.GetWindowAttributes(wm.xUtil.Conn(), xWinID)

	windowAttr, err := windowAttrReq.Reply()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to retrieve attributes of window '%#x", xWinID)

		return
	}

	isDisplayable := (windowAttr.MapState == xproto.MapStateUnmapped || windowAttr.MapState == xproto.MapStateViewable)

	res = (winState != nil && winTypeIsManageable && isDisplayable)

	return res, err
}

// func (wm *WindowManager) getManageableWindows(xWinID xproto.Window) (res *treeset.Set, err error) {
// 	res = treeset.NewWith(WindowIDComparator)

// 	isManageable, manageableStateErr := wm.isWindowManageable(xWinID)
// 	if manageableStateErr == nil && isManageable {
// 		internal.VerboseLogger.Debugf("Window '%#x' is manageable", xWinID)
// 		res.Add(xWinID)
// 	} else {
// 		//internal.VerboseLogger.Debugf("Window '%#x' is not manageable", xWinID)
// 	}

// 	qTreeReq := xproto.QueryTree(wm.xUtil.Conn(), xWinID)

// 	qTree, err := qTreeReq.Reply()
// 	if err != nil {
// 		internal.NormalLogger.Warnf("Unable to retrieve xwindows tree: %s", err)

// 		return
// 	}

// 	for _, child := range qTree.Children {
// 		grandChildren, childrenErr := wm.getManageableWindows(child)
// 		if childrenErr != nil {
// 			continue
// 		}
// 		res.Add(grandChildren.Values()...)
// 	}

// 	return
// }

func (wm *WindowManager) loadManageableWindows(xWinID xproto.Window, indent string) (err error) {
	//internal.VerboseLogger.Debugf("%s Processing '%#x'", indent, xWinID)
	isManageable, manageableStateErr := wm.isWindowManageable(xWinID)
	if manageableStateErr == nil && isManageable {
		internal.VerboseLogger.Debugf("Window '%#x' is manageable", xWinID)
		var weh *WindowEventHandler
		weh, err = GetWindowEventHandlerInstance()
		if err != nil {
			err = eris.Wrap(err, "Unable to get the instance of the window event handler")

			return
		}
		if handleErr := weh.handleWindowCreation(xWinID); handleErr != nil {
			internal.VerboseLogger.Warnf("Unable to handle creation of window '%#x'", xWinID)
			err = eris.Wrapf(handleErr, "Unable to handle creation of window '%#x'", xWinID)
		}
		return
	} else {
		//internal.VerboseLogger.Debugf("Window '%#x' is not manageable", xWinID)
	}

	qTreeReq := xproto.QueryTree(wm.xUtil.Conn(), xWinID)

	qTree, err := qTreeReq.Reply()
	if err != nil {
		internal.NormalLogger.Warnf("Unable to retrieve xwindows tree: %s", err)

		return
	}

	for _, child := range qTree.Children {
		childrenErr := wm.loadManageableWindows(child, indent+"  ")
		if childrenErr != nil {
			continue
		}
	}

	return
}

func (wm *WindowManager) LoadWindows() (err error) {

	err = wm.loadManageableWindows(wm.xUtil.RootWin(), "")
	if err != nil {
		err = eris.Wrap(internal.InternalError, "Unable to get list of current windows")

		return
	}

	if wm.windows.Size() == 0 {
		internal.VerboseLogger.Warnf("There are currently no windows to manage")
	}

	return err
}

// func (wm *WindowManager) LoadWindows() (err error) {
// 	var weh *WindowEventHandler

// 	weh, err = GetWindowEventHandlerInstance()
// 	if err != nil {
// 		err = eris.Wrap(err, "Unable to get the instance of the window event handler")

// 		return
// 	}

// 	windows, err := wm.getManageableWindows(wm.xUtil.RootWin())
// 	if err != nil {
// 		err = eris.Wrap(internal.InternalError, "Unable to get list of current windows")

// 		return
// 	}

// 	if windows.Size() == 0 {
// 		internal.VerboseLogger.Warnf("There are currently no windows to manage")
// 	}

// 	for _, xWinID := range windows.Values() {
// 		if handleErr := weh.handleWindowCreation(xWinID.(xproto.Window)); handleErr != nil {
// 			internal.VerboseLogger.Warnf("Unable to handle creation of window '%#x'", xWinID)
// 			err = eris.Wrapf(handleErr, "Unable to handle creation of window '%#x'", xWinID)
// 		}
// 	}

// 	return err
// }

func (wm *WindowManager) GetFocusedWindow() (res *Window, err error) {
	focusedXWinID, err := ewmh.ActiveWindowGet(wm.xUtil)
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to obtain the active window")

		return
	}

	internal.NormalLogger.Debugf("Window '%#x' seems to be active", focusedXWinID)

	focusedWindowSet := wm.GetWindows(func(win *Window) bool {
		return win.XWinID() == focusedXWinID
	})

	if focusedWindowSet.Size() != 1 {
		err = eris.Wrapf(internal.NotExistsError, "Window '%#x' is focused but is not managed by snapi3", focusedXWinID)

		return
	}

	res = focusedWindowSet.Values()[0].(*Window)

	return
}

func (wm *WindowManager) getNodes(parentNode *i3.Node, predicate func(*i3.Node) bool) (res *treeset.Set, err error) {
	if parentNode == nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid parent node")

		return
	}

	res = treeset.NewWith(NodeComparator)

	for _, n := range parentNode.Nodes {
		if n != nil && predicate(n) {
			res.Add(n)
		} else {
			var children *treeset.Set

			children, err = wm.getNodes(n, predicate)
			if err != nil {
				err = eris.Wrap(internal.InternalError, "Unable to retrieve non-floating nodes matching the predicate")

				return
			}

			res.Add(children.Values()...)
		}
	}

	for _, n := range parentNode.FloatingNodes {
		if n != nil && predicate(n) {
			res.Add(n)
		}

		var children *treeset.Set

		children, err = wm.getNodes(n, predicate)
		if err != nil {
			err = eris.Wrapf(err, "Unable to retrieve floating nodes matching the predicate")

			return
		}

		res.Add(children.Values()...)
	}

	return res, err
}

func (wm *WindowManager) GetNode(nodeID i3.NodeID) (res *i3.Node, err error) {
	nodeSet, err := wm.GetNodes(func(node *i3.Node) bool {
		return node.ID == nodeID
	})
	if err != nil {
		err = eris.Wrapf(err, "Unable to find the node with ID '%d'", nodeID)

		return
	}

	if nodeSet.Size() != 1 {
		err = eris.Wrapf(internal.NotExistsError, "Unable to find the node with ID '%d'", nodeID)

		return
	}

	res = nodeSet.Values()[0].(*i3.Node)

	return
}

func (wm *WindowManager) GetNodes(predicate func(*i3.Node) bool) (res *treeset.Set, err error) {
	i3Tree, err := i3.GetTree()
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to retrieve the i3 tree")

		return
	}

	focusedOutput := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Type == i3.OutputNode
	})
	if focusedOutput == nil {
		err = eris.Wrapf(internal.InternalError, "Unable to retrieve focused output")

		return
	}

	res, err = wm.getNodes(i3Tree.Root, predicate)

	return
}
