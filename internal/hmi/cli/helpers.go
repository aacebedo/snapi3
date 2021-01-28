package cli

import (
	"github.com/aacebedo/snapi3/internal/hmi"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/rotisserie/eris"
)

func FilterWindowsToProcess(winsToFilter *treeset.Set, filterStrs []string) (res *treeset.Set, err error) {
	res = treeset.NewWith(windowmgt.WindowComparator)
	filters := doublylinkedlist.New()
	for _, filterStr := range filterStrs {
		filter, filterCreationErr := windowmgt.NewWindowFilterWithStr(filterStr)
		if filterCreationErr == nil {
			filters.Add(filter)
		}
	}

	for winIt := winsToFilter.Iterator(); winIt.Next(); {
		curWin := winIt.Value().(*windowmgt.Window)
		shallBeSnapped := filters.Size() != 0
		for filterIt := filters.Iterator(); filterIt.Next(); {
			shallBeSnapped = shallBeSnapped && filterIt.Value().(*windowmgt.WindowFilter).IsMatching(curWin)
		}
		if shallBeSnapped {
			res.Add(curWin)

		}
	}
	return
}

func GetWindowsToProcess(xWinIDStr *string, winsToFilter *treeset.Set, filterStrs []string) (res *treeset.Set, err error) {
	res = treeset.NewWith(windowmgt.WindowComparator)
	if len(filterStrs) == 0 {
		var targetedWin *windowmgt.Window
		targetedWin, err = hmi.GetTargetedWindow(xWinIDStr)
		if err != nil {
			err = eris.Wrapf(err, "Unable to obtain the targeted window")
			return
		}
		res.Add(targetedWin)
	} else {
		res, err = FilterWindowsToProcess(winsToFilter, filterStrs)
		if err != nil {
			err = eris.Wrapf(err, "Unable to snap the windows: cannot filter windows")

			return
		}
	}
	return
}
