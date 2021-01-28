package windowmgt

import (
	"sync"

	"github.com/aacebedo/snapi3/internal"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"github.com/rotisserie/eris"
)

type WindowGroupIdentifier uint8

type WindowGroup struct {
	filters  *doublylinkedlist.List
	windows  *treeset.Set
	name     string
	color    Color
	winMutex sync.Mutex
	id       WindowGroupIdentifier
}

func WindowGroupComparator(left, right interface{}) int {
	leftGroup := left.(*WindowGroup)
	rightGroup := right.(*WindowGroup)

	return utils.UInt8Comparator(uint8(leftGroup.ID()), uint8(rightGroup.ID()))
}

func NewWindowGroup(id WindowGroupIdentifier, name string, color Color) (res *WindowGroup) {
	res = &WindowGroup{id: id, name: name, color: color, windows: treeset.NewWith(WindowComparator), filters: doublylinkedlist.New()}

	return
}

func (wg *WindowGroup) ID() (res WindowGroupIdentifier) {
	res = wg.id

	return
}

func (wg *WindowGroup) Name() (res string) {
	res = wg.name

	return
}

func (wg *WindowGroup) SetName(name string) {
	wg.name = name
}

func (wg *WindowGroup) Color() (res Color) {
	res = wg.color

	return
}

func (wg *WindowGroup) Filters() (res *doublylinkedlist.List) {
	res = wg.filters

	return
}

func (wg *WindowGroup) SetColor(color Color) {
	wg.color = color
}

func (wg *WindowGroup) AddFilter(winFilter *WindowFilter) {
	wg.filters.Add(winFilter)
}

func (wg *WindowGroup) AddWindow(winToAdd *Window) (err error) {
	wg.winMutex.Lock()
	defer wg.winMutex.Unlock()

	if wg.windows.Contains(winToAdd) {
		err = eris.Wrapf(internal.AlreadyExistsError, "Group '%s' already contains window '%s'", wg.Name(), winToAdd.Name())

		return
	}

	wg.windows.Add(winToAdd)
	internal.VerboseLogger.Debugf("Added window '%s' to group '%s'", winToAdd.Name(), wg.Name())

	return
}

func (wg *WindowGroup) RemoveWindow(winToAdd *Window) (err error) {
	wg.winMutex.Lock()
	defer wg.winMutex.Unlock()

	if !wg.windows.Contains(winToAdd) {
		err = eris.Wrapf(internal.NotExistsError, "Group '%s' does not contain window '%s'", wg.Name(), winToAdd.Name())

		return
	}

	internal.VerboseLogger.Debugf("Removed window '%s' from group '%s'", winToAdd.Name(), wg.Name())
	wg.windows.Remove(winToAdd)

	return
}

func (wg *WindowGroup) GetWindows() (res *treeset.Set) {
	res = wg.windows

	return
}

func (wg *WindowGroup) WindowIsMatchingFilters(winToFilter *Window) (res bool) {
	res = false

	if wg.filters.Size() > 0 {
		filterIt := wg.filters.Iterator()
		filter := filterIt.Value().(WindowFilter)
		res = filter.IsMatching(winToFilter)

		for filterIt.Next(); filterIt.Next(); {
			filter = filterIt.Value().(WindowFilter)
			switch filter.Operator() {
			case internal.Or:
				res = res || filter.IsMatching(winToFilter)
			case internal.And:
				res = res && filter.IsMatching(winToFilter)
			}
		}
	}

	return
}

func (wg *WindowGroup) Contains(winToSearch *Window) (res bool) {
	wg.winMutex.Lock()
	defer wg.winMutex.Unlock()

	return wg.windows.Contains(winToSearch)
}
