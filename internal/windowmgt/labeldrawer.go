package windowmgt

import (
	"math"
	"sync"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/aacebedo/snapi3/internal"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/utils"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

const (
	PillXOffset uint32 = 5
	PillYOffset uint32 = 5
)

var (
	ldInstance     *LabelDrawer    //nolint:gochecknoglobals //Needed for the singleton pattern
	ldInstanceLock = &sync.Mutex{} //nolint:gochecknoglobals //Needed for the singleton pattern
)

type LabelDrawer struct {
	xgbConn            *xgb.Conn
	labelWins          *hashmap.Map
	labelWinsDisplayed bool
	winMutex           sync.Mutex
}

func PillWidthComparator(left, right interface{}) int {
	leftPill := left.(*Pill)
	rightPill := right.(*Pill)
	leftPw, _ := leftPill.GetDimensions()
	rightPw, _ := rightPill.GetDimensions()

	return utils.UInt32Comparator(leftPw, rightPw)
}

func GetLabelDrawerInstance() (res *LabelDrawer, err error) {
	if ldInstance == nil {
		ldInstanceLock.Lock()
		defer ldInstanceLock.Unlock()

		if ldInstance == nil {
			var xgbConn *xgb.Conn

			xgbConn, err = xgb.NewConn()
			if err != nil {
				err = eris.Wrap(internal.InternalError, "Unable to open an xgb connection")

				return
			}

			ldInstance = &LabelDrawer{labelWins: hashmap.New(), xgbConn: xgbConn, labelWinsDisplayed: false}
		}
	}

	res = ldInstance

	return
}

func (ld *LabelDrawer) Start() {
	gtk.Init(nil)

	go gtk.Main()
	internal.NormalLogger.Debugf("Label drawer is started")
}

func (ld *LabelDrawer) Stop() {
	//TODO change the time
	time.Sleep(2 * time.Second)
	gtk.MainQuit()
	internal.NormalLogger.Debugf("Label drawer is stopped")
}

func (ld *LabelDrawer) GetLabelFontExtents(cr *cairo.Context) (res cairo.FontExtents) {
	cr.Save()
	cr.SelectFontFace(viper.GetString("group_labels.font"), cairo.FONT_SLANT_NORMAL, cairo.FONT_WEIGHT_BOLD)
	cr.SetFontSize(viper.GetFloat64("group_labels.font_size"))
	res = cr.FontExtents()
	cr.Restore()

	return
}

// func (ld *LabelDrawer) HideAllWindowsLabels() {
// 	ld.winMutex.Lock()
// 	defer ld.winMutex.Unlock()

// 	if ld.labelWinsDisplayed {

// 		ld.labelWinsDisplayed = false
// 	}
// }

// func (ld *LabelDrawer) ShowAllWindowLabels() (err error) {
// 	ld.winMutex.Lock()
// 	defer ld.winMutex.Unlock()
// 	if !ld.labelWinsDisplayed {
// 		ld.labelWinsDisplayed = true
// 	}

// 	wins := GetWindowManagerInstance().GetWindows(func(win *Window) bool { return true })

// 	for winIt := wins.Iterator(); winIt.Next(); {
// 		curWin := winIt.Value().(*Window)

// 		labelErr := ld.showWindowLabels(curWin)
// 		if labelErr != nil {
// 			err = eris.Wrapf(err, "Unable to display window labels of window '%s'", curWin.Name())

// 			return
// 		}
// 	}
// 	return
//}

func (ld *LabelDrawer) drawLabels(win *Window, labelWin *gtk.Window, cr *cairo.Context) {
	wgm := GetWindowGroupManagerInstance()
	groups := wgm.GetGroupsOfWindow(win)

	parenWinGeom, _ := xproto.GetGeometry(ld.xgbConn, xproto.Drawable(win.XWinID())).Reply()
	ww := uint32(parenWinGeom.Width)
	wh := uint32(0)
	curX := PillXOffset
	curY := PillYOffset
	_, lineHeight := NewPill("Dummy", *NewColor(0, 0, 0), cr).GetDimensions()

	if groups.Size() != 0 {
		widthOrderedPills := doublylinkedlist.New()

		groupIt := groups.Iterator()
		for groupIt.Next() {
			curGroup := groupIt.Value().(*WindowGroup)
			p := NewPill(curGroup.Name(), curGroup.Color(), cr)
			widthOrderedPills.Add(p)
		}
		widthOrderedPills.Sort(PillWidthComparator)

		groupIdx := 1
		linesNb := uint32(1)

		widthOrderedPillsIt := widthOrderedPills.Iterator()
		for widthOrderedPillsIt.End(); widthOrderedPillsIt.Prev(); {
			curPill := widthOrderedPillsIt.Value().(*Pill)

			pw, _ := curPill.GetDimensions()
			if groupIdx != 1 && curX+pw+PillXOffset > ww-PillXOffset*2 {
				curY += lineHeight + PillYOffset
				curX = PillXOffset
				linesNb++
			}

			curPill.draw(float64(curX), float64(curY))
			curX += pw + PillXOffset
			groupIdx++
		}

		wh = linesNb*lineHeight + PillYOffset*linesNb + 2*PillYOffset
	} else {
		p := NewPill("<no groups>", *NewColor(math.MaxUint8, math.MaxUint8, math.MaxUint8), cr)

		p.draw(float64(curX), float64(curY))
		wh = lineHeight + 2*PillYOffset
	}

	labelWin.SetSizeRequest(int(ww), int(wh))
}

func (ld *LabelDrawer) ShowWindowLabels(win *Window) (err error) {
	if win == nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Window argument is nil")

		return
	}

	ld.winMutex.Lock()
	defer ld.winMutex.Unlock()

	internal.NormalLogger.Debugf("Showing labels for window  '%s'", win.Name())

	labelWinItf, labelWinExists := ld.labelWins.Get(win.XWinID())
	if !labelWinExists {
		_, err = glib.IdleAdd(func() (res bool) {
			ld.winMutex.Lock()
			defer ld.winMutex.Unlock()
			res = false
			labelWin, _ := gtk.WindowNew(gtk.WINDOW_POPUP)
			ld.labelWins.Put(win.XWinID(), labelWin)
			da, _ := gtk.DrawingAreaNew()
			labelWin.Add(da)
			labelWin.SetFocusOnMap(false)
			labelWin.SetResizable(false)
			labelWin.SetDefaultSize(0, 0)
			labelWin.SetSizeRequest(0, 0)
			labelWin.SetDestroyWithParent(true)
			labelWin.SetDecorated(false)

			_, err = da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
				ld.drawLabels(win, labelWin, cr)
			})
			if err != nil {
				return
			}
			labelWin.ShowAll()
			gdkWin, _ := labelWin.GetWindow()

			if err = xproto.ReparentWindowChecked(ld.xgbConn, xproto.Window(gdkWin.GetXID()), win.XWinID(), 0, 0).Check(); err != nil {
				err = eris.Wrapf(internal.InternalError, "Unable to reparent window '%s'", win)

				return
			}

			return
		})
		if err != nil {
			err = eris.Wrap(internal.InternalError, "Unable to draw window group labels with gtk commands")

			return
		}
	} else {
		labelWin, _ := labelWinItf.(*gtk.Window)
		labelWin.QueueDraw()
	}

	return err
}

func (ld *LabelDrawer) HideWindowLabels(win *Window) (err error) {
	if win == nil {
		err = eris.Wrapf(internal.InvalidArgumentError, "Window argument is nil")

		return
	}

	ld.winMutex.Lock()
	defer ld.winMutex.Unlock()

	internal.NormalLogger.Debugf("Hiding labels for window '%s'", win.Name())

	ld.labelWins.Remove(win.XWinID())

	return err
}
