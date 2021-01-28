package windowmgt

import (
	"sync"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/aacebedo/snapi3/internal"
	"github.com/rotisserie/eris"
	"go.i3wm.org/i3/v4"
)

type WindowEventHandler struct {
	xutil *xgbutil.XUtil
}

var (
	wehInstance     *WindowEventHandler //nolint:gochecknoglobals //Needed for singleton pattern
	wehInstanceLock = &sync.Mutex{}     //nolint:gochecknoglobals //Needed for singleton pattern
)

func GetWindowEventHandlerInstance() (res *WindowEventHandler, err error) {
	if wehInstance == nil {
		wehInstanceLock.Lock()
		defer wehInstanceLock.Unlock()

		if wehInstance == nil {
			var xutil *xgbutil.XUtil

			xutil, err = xgbutil.NewConn()
			if err != nil {
				err = eris.Wrap(internal.InternalError, "Unable to create a xgb connection")

				return
			}

			wehInstance = &WindowEventHandler{xutil: xutil}
		}
	}

	res = wehInstance

	return
}

func (weh *WindowEventHandler) StartEventProcessing() {
	go weh.processX11Events()
	go weh.processI3Events()
}

func (weh *WindowEventHandler) processX11Events() {
	xevent.Main(weh.xutil)
}

func (weh *WindowEventHandler) processI3Events() {
	er := i3.Subscribe(i3.WindowEventType)
	for er.Next() {
		ev := er.Event().(*i3.WindowEvent)
		if ev == nil {
			internal.NormalLogger.Error("Unable to cast to WindowEvent type")

			continue
		}

		switch ev.Change {
		case "new":
			if err := weh.handleWindowCreation(xproto.Window(ev.Container.Window)); err != nil {
				internal.NormalLogger.Warnf("Unable to handle the creation of window '%#x'", ev.Container.Window)
			}
		case "close":
			if err := weh.handleWindowDestruction(xproto.Window(ev.Container.Window)); err != nil {
				internal.NormalLogger.Warnf("Unable to handle the destruction of window '%#x'", ev.Container.Window)
			}
		default:
		}
	}
}

func (weh *WindowEventHandler) handleWindowCreation(xWinID xproto.Window) (err error) {
	internal.NormalLogger.Debugf("Handling creation of the window '%#x'", xWinID)
	isWindowManageable, err := wmInstance.isWindowManageable(xWinID)

	if err == nil && isWindowManageable {
		var win *Window

		var wm *WindowManager

		wm, err = GetWindowManagerInstance()
		if err != nil {
			err = eris.Wrapf(err, "Unable to obtain the window manager instance")

			return
		}

		win, err = wm.AddWindow(xWinID)
		if err != nil {
			err = eris.Wrapf(err, "Unable to manage newly created window '%#x'", xWinID)

			return
		}

		xwin := xwindow.New(weh.xutil, win.XWinID())
		if err = xwin.Listen(xproto.EventMaskStructureNotify); err != nil {
			err = eris.Wrap(err, "Unable to listen to substructure notification of root window")

			return
		}

		xevent.ConfigureNotifyFun(
			func(X *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
				var ld *LabelDrawer
				ld, err = GetLabelDrawerInstance()

				if err != nil {
					internal.NormalLogger.Warnf("Unable to retrieve the instance of the label drawer: %s", err)

					return
				}

				if err = ld.ShowWindowLabels(win); err != nil {
					internal.NormalLogger.Warnf("Unable to show the labels for win '%s': %s", win, err)

					return
				}
			}).Connect(weh.xutil, xWinID)

		var ld *LabelDrawer

		ld, err = GetLabelDrawerInstance()
		if err != nil {
			err = eris.Wrap(err, "Unable to retrieve the instance of the label drawer")

			return
		}

		if err = ld.ShowWindowLabels(win); err != nil {
			err = eris.Wrapf(err, "Unable to show the labels for win '%s'", win)

			return
		}
	} else {
		internal.NormalLogger.Warnf("Unable to handle window '%#x'", xWinID)
	}

	return err
}

func (weh *WindowEventHandler) handleWindowDestruction(winID xproto.Window) (err error) {
	var wm *WindowManager

	wm, err = GetWindowManagerInstance()
	if err != nil {
		err = eris.Wrapf(err, "Unable to obtain the window manager instance")

		return
	}

	win, err := wm.GetWindow(winID)
	if err != nil {
		err = eris.Wrapf(err, "Unable find window '%#x'", winID)

		return
	}

	ld, err := GetLabelDrawerInstance()
	if err != nil {
		err = eris.Wrap(err, "Unable to retrieve the instance of the label drawer")

		return
	}

	if err = ld.HideWindowLabels(win); err != nil {
		err = eris.Wrap(err, "Unable to hide the window labels")

		return
	}

	if err = wm.RemoveWindow(win.XWinID()); err != nil {
		err = eris.Wrapf(err, "Unable to remove the window '%#x' from the window manager", win.XWinID())

		return
	}

	return err
}
