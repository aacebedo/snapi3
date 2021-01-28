package windowmgt

import (
	"crypto/rand"
	"sync"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/aacebedo/snapi3/internal"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

// type windowGroupDescriptor struct {
// 	ID    WindowGroupIdentifier `yaml:"id"`
// 	Name  string                `yaml:"name"`
// 	Color string                `yaml:"color"`
// 	//Windows []uint32 `yaml:"windows"`
// }

// type windowGroupsDescriptor struct {
// 	Groups []windowGroupDescriptor `yaml:"groups"`
// }

var (
	wgmInstance     *WindowGroupManager //nolint:gochecknoglobals //Needed for singleton pattern
	wgmInstanceLock = &sync.Mutex{}     //nolint:gochecknoglobals //Needed for singleton pattern
)

type WindowGroupManager struct {
	groups   *treeset.Set
	winMutex sync.Mutex
}

func GetWindowGroupManagerInstance() (res *WindowGroupManager) {
	if wgmInstance == nil {
		wgmInstanceLock.Lock()
		defer wgmInstanceLock.Unlock()

		if wgmInstance == nil {
			wgmInstance = &WindowGroupManager{groups: treeset.NewWith(WindowGroupComparator)}
		}
	}

	res = wgmInstance

	return
}

func (wgm *WindowGroupManager) GetGroupsOfWindow(win *Window) (res *treeset.Set) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()
	res = wgm.groups.Select(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return curGroup.Contains(win)
	})

	return
}

func (wgm *WindowGroupManager) GetGroups(predicate func(*WindowGroup) bool) (res *treeset.Set) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()
	res = wgm.groups.Select(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return predicate(curGroup)
	})

	return
}

func (wgm *WindowGroupManager) GetGroup(groupID WindowGroupIdentifier) (res *WindowGroup, err error) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()
	_, groupItf := wgm.groups.Find(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return curGroup.ID() == groupID
	})

	if groupItf == nil {
		err = eris.Wrapf(internal.NotExistsError, "Cannot retrieve non existent group '%d'", groupID)

		return
	}

	res = groupItf.(*WindowGroup)

	return
}

func (wgm *WindowGroupManager) Contains(groupID WindowGroupIdentifier) (res bool) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()
	_, groupItf := wgm.groups.Find(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return curGroup.ID() == groupID
	})

	res = groupItf != nil

	return
}
func (wgm *WindowGroupManager) GenerateGroupID() (res WindowGroupIdentifier, err error) {
	res = 0

	values := wgm.groups.Values()
	if len(values) == 256 {
		err = eris.Wrapf(internal.OutOfRangeArgumentError, "Cannot add more than 256 groups")
		return
	}
	idx := 0
	for ; idx < len(values)-1; idx++ {
		if values[idx+1].(*WindowGroup).ID()-values[idx].(*WindowGroup).ID() > 1 {
			res = values[idx].(*WindowGroup).ID() + 1

			break
		}
	}
	if idx == len(values)-1 {
		res = values[idx].(*WindowGroup).ID() + 1
	}
	return
}

func (wgm *WindowGroupManager) AddGroup(groupID WindowGroupIdentifier, groupName string, color Color) (res *WindowGroup, err error) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()

	_, groupItf := wgm.groups.Find(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return curGroup.ID() == groupID
	})

	if groupItf != nil {
		err = eris.Wrapf(internal.AlreadyExistsError, "Group with ID '%d' already exists", groupID)

		return
	}

	res = NewWindowGroup(groupID, groupName, color)
	wgm.groups.Add(res)
	internal.VerboseLogger.Debugf("Added group '%s'", res.Name())

	return
}

func (wgm *WindowGroupManager) RemoveGroup(groupIDToRemove WindowGroupIdentifier) (err error) {
	wgm.winMutex.Lock()
	defer wgm.winMutex.Unlock()

	_, groupItf := wgm.groups.Find(func(index int, value interface{}) bool {
		curGroup := value.(*WindowGroup)

		return curGroup.ID() == groupIDToRemove
	})

	if groupItf == nil {
		err = eris.Wrapf(internal.AlreadyExistsError, "Group with ID '%d' does not exist", groupIDToRemove)

		return
	}

	groupToRemove := groupItf.(*WindowGroup)
	wgm.groups.Remove(groupToRemove)
	internal.VerboseLogger.Debugf("Removed group '%s'", groupToRemove.Name())

	return
}

func (wgm *WindowGroupManager) LoadWindowGroups() (err error) {

	var config internal.Configuration

	if err = viper.Unmarshal(&config); err != nil {
		internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

		return
	}

	/*internal.VerboseLogger.Debugf("Loading window groups from file '%s'", windowGroupsFilepath)
	winGroupsDesc := &windowGroupsDescriptor{}

	file, err := os.Open(windowGroupsFilepath)
	if err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to open window groups description file '%s'", windowGroupsFilepath)

		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&winGroupsDesc); err != nil {
		err = eris.Wrapf(internal.InternalError, "Unable to parse YAML string contained in '%s'", windowGroupsFilepath)

		return
	}*/

	wgm.winMutex.Lock()
	wgm.groups.Clear()
	wgm.winMutex.Unlock()

	for _, groupConfig := range config.GroupConfigs {
		var groupColor *Color

		groupColor, err = NewColorFromHTMLCode(groupConfig.Color)
		if err != nil {
			internal.VerboseLogger.Warnf("Group '%s' ignored because its color '%s' is invalid", groupConfig.Name, groupConfig.Color)

			buf := make([]byte, 1)

			_, _ = rand.Read(buf)
			rComponent := buf[0]

			_, _ = rand.Read(buf)
			gComponent := buf[0]

			_, _ = rand.Read(buf)
			bComponent := buf[0]

			groupColor = NewColor(rComponent, gComponent, bComponent)
			internal.VerboseLogger.Warnf("Replacing color of group '%s' with '%s'", groupConfig.Name, *groupColor)

			continue
		}

		//var wg *WindowGroup

		newGroup, err := wgm.AddGroup(WindowGroupIdentifier(groupConfig.ID), groupConfig.Name, *groupColor)
		if err != nil {
			internal.VerboseLogger.Warnf("Group '%s' already defined", groupConfig.Name)

			continue
		}
		for _, filterConfig := range groupConfig.Filters {
			if len(filterConfig.Filters) == 0 {
				filter, filterCreateErr := NewWindowFilter(filterConfig.WinProperty, filterConfig.Regex, filterConfig.Operator)
				if filterCreateErr == nil {
					newGroup.AddFilter(filter)
				}
			}
		}

		// for _, xWinID := range winGroupDesc.Windows {
		// 	var win *Window

		// 	if err == nil {

		// 		if err = wg.AddWindow(win); err != nil {
		// 			internal.VerboseLogger.Warnf("Window '%s' already belongs the window group '%s'", win.Name(), winGroupDesc.Name)
		// 		} else {
		// 			internal.VerboseLogger.Debugf("Window '%s' added group '%s'", win.Name(), winGroupDesc.Name)
		// 		}
		// 	} else {
		// 		internal.VerboseLogger.Infof("Window '%0x' in group '%s' does not exist anymore, it won't be handled", xWinID, winGroupDesc.Name)
		// 	}
		// }
	}

	wm, err := GetWindowManagerInstance()
	if err != nil {
		err = eris.Wrapf(err, "Unable to retrieved the window manager instance")

		return
	}

	for winIt := wm.GetWindows(func(window *Window) bool { return true }).Iterator(); winIt.Next(); {
		var groupIDs []uint8

		win := winIt.Value().(*Window)
		groupIDs, _, err = win.GetAtom("SNAPI3_GROUPS")
		for _, groupID := range groupIDs {
			group, err := wgm.GetGroup(WindowGroupIdentifier(groupID))
			if err != nil {
				internal.NormalLogger.Warnf("Window '%s' belonged to the inexistent group '%d'", win.Name(), groupID)
			}
			group.AddWindow(win)
		}

		groups := wgm.GetGroups(func(*WindowGroup) bool { return true })
		for groupIt := groups.Iterator(); groupIt.Next(); {
			curGroup := groupIt.Value().(*WindowGroup)
			if curGroup.WindowIsMatchingFilters(win) {
				curGroup.AddWindow(win)
			}
		}
	}
	return err
}

func (wgm *WindowGroupManager) SaveWindowGroups() (err error) {
	// windowGroupsDirpath := filepath.Dir(windowGroupsFilepath)

	// if _, err = os.Stat(windowGroupsDirpath); os.IsNotExist(err) {
	// 	err = os.MkdirAll(windowGroupsDirpath, os.ModePerm)
	// }

	// if err != nil {
	// 	err = eris.Wrapf(internal.InternalError, "Unable to create window groups variable file directory '%s'", windowGroupsDirpath)

	// 	return
	// }

	// winGroupsDesc := &windowGroupsDescriptor{}
	// var config internal.Configuration

	// if err = viper.Unmarshal(&config); err != nil {
	// 	internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
	// 	err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

	// 	return
	// }
	{
		wgm.winMutex.Lock()

		groupConfigs := make([]internal.GroupConfiguration, 0)
		for groupIt := wgm.groups.Iterator(); groupIt.Next(); {
			groupConfig := internal.GroupConfiguration{}

			curGroup := groupIt.Value().(*WindowGroup)
			groupConfig.Name = curGroup.Name()
			groupColor := curGroup.Color()
			groupConfig.Color = groupColor.ToHTMLCode()
			groupConfig.ID = uint8(curGroup.ID())
			for filterIt := curGroup.Filters().Iterator(); filterIt.Next(); {
				groupConfig.Filters = append(groupConfig.Filters, filterIt.Value().(*WindowFilter).ConvertToConfig())
			}
			groupConfigs = append(groupConfigs, groupConfig)

			//curGroupDesc := windowGroupDescriptor{Name: curGroup.Name(), Color: groupColor.ToHTMLCode()}
			// for winIt := curGroup.GetWindows().Iterator(); winIt.Next(); {

			// }

			//curGroupDesc.Windows = append(curGroupDesc.Windows, uint32(winIt.Value().(*Window).XWinID()))

			//		winGroupsDesc.Groups = append(winGroupsDesc.Groups, curGroupDesc)
		}
		wgm.winMutex.Unlock()

		viper.Set("groups", groupConfigs)
		viper.WriteConfig()
	}
	wm, err := GetWindowManagerInstance()
	//TODO: Check error
	for winIt := wm.GetWindows(func(window *Window) bool { return true }).Iterator(); winIt.Next(); {
		var IDs []uint8

		winGroups := wgm.GetGroupsOfWindow(winIt.Value().(*Window))

		for groupIt := winGroups.Iterator(); groupIt.Next(); {
			IDs = append(IDs, uint8(groupIt.Value().(*WindowGroup).ID()))
		}
		err = winIt.Value().(*Window).SetAtom(xproto.AtomInteger, "SNAPI3_GROUPS", 8, []byte(IDs), uint32(len(IDs)))
		//TODO : check errors
	}

	/*	file, err := os.OpenFile(windowGroupsFilepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			err = eris.Wrapf(internal.InternalError, "Unable to open window groups description file '%s'", windowGroupsFilepath)

			return
		}
		defer file.Close()

		encoder := yaml.NewEncoder(file)

		if err = encoder.Encode(&winGroupsDesc); err != nil {
			err = eris.Wrapf(internal.InternalError, "Unable to save YAML string contained in '%s'", windowGroupsFilepath)

			return
		}*/

	return err
}
