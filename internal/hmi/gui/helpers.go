package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/aacebedo/gorofimenus/pkg/gorofimenus"
	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/alessio/shellescape"
	clicmd "github.com/commander-cli/cmd"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasttemplate"
)

func positionToAnchorAndLocation(position string) (location, anchor string, err error) {
	switch position {
	case "left-top":
		anchor = "north west" //nolint:goconst //No reason to create a const from that
		location = "north west"
	case "top":
		anchor = "north" //nolint:goconst //No reason to create a const from that
		location = "north"
	case "right-top":
		anchor = "north east" //nolint:goconst //No reason to create a const from that
		location = "north east"
	case "left":
		anchor = "west" //nolint:goconst //No reason to create a const from that
		location = "west"
	case "center": //nolint:goconst //No reason to create a const from that
		anchor = "center"
		location = "center"
	case "right":
		anchor = "east" //nolint:goconst //No reason to create a const from that
		location = "east"
	case "bottom":
		anchor = "south west" //nolint:goconst //No reason to create a const from that
		location = "south west"
	case "center-bottom":
		anchor = "south" //nolint:goconst //No reason to create a const from that
		location = "south"
	case "right-bottom":
		anchor = "south east" //nolint:goconst //No reason to create a const from that
		location = "south east"
	default:
		err = eris.Wrapf(internal.InvalidArgumentError, "Invalid value of position '%s'", position)

		return
	}

	return location, anchor, err
}

func askUser(themeOptionalPartFilepath, themeMandatoryPart string,
	themeMandatoryPartElements map[string]interface{},
	menuItems *linkedhashmap.Map, menuName string) (selectedItemValue string, selectedItemPayload interface{}, err error) {
	tmpThemeFile, err := generateThemeFile(themeOptionalPartFilepath, themeMandatoryPart, themeMandatoryPartElements)
	if err != nil {
		err = eris.Wrapf(err, "Unable to generate theme file for cell selection")

		return
	}

	defer os.Remove(tmpThemeFile)

	gorofimenus.SetLogLevel(logrus.DebugLevel, 2)

	menu := gorofimenus.NewMenu(menuName, nil)
	for it := menuItems.Iterator(); it.Next(); {
		name := it.Key().(string)

		if err = menu.AddSubMenu(gorofimenus.NewMenu(name, it.Value())); err != nil {
			err = eris.Wrap(internal.InternalError, "Unable to build the selection grid")

			return
		}
	}

	menu.SetOptions([]string{"-no-custom", "-theme", tmpThemeFile})

	selection, err := menu.GetSelection(true)
	if err != nil {
		internal.VerboseLogger.Errorf("Unable to get the selected cell from the user: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to get the selected cell from the user")

		return
	}

	selectionIt := selection.Iterator()
	selectionIt.Last()
	selectedMenu := selectionIt.Value().(*gorofimenus.Menu)

	selectedItemValue = selectedMenu.Value()
	selectedItemPayload = selectedMenu.Payload()

	return selectedItemValue, selectedItemPayload, err
}

func generateThemeFile(themeOptionalPartFilepath, themeMandatoryPart string,
	themeMandatoryPartElements map[string]interface{}) (res string, err error) {
	tmpThemeFile, err := ioutil.TempFile(os.TempDir(), "snapi3-theme")
	if err != nil {
		err = eris.Wrap(internal.InternalError, "Unable to create the temporary file for snapi3 grid theme")

		return
	}

	themeMandatoryPartTemplate := fasttemplate.New(themeMandatoryPart, "{{", "}}")
	instantiatedThemeMandatoryPart := themeMandatoryPartTemplate.ExecuteString(themeMandatoryPartElements)

	themeOptionalPartFile, err := os.Open(themeOptionalPartFilepath)
	if err != nil {
		//TODO:deal with the error
		fmt.Println(err)

		return
	}

	themeOptionalPart, err := ioutil.ReadAll(themeOptionalPartFile)

	themeFullTemplateStr :=
		`{{theme_optional_part}}
{{theme_mandatory_part}}`

	themeFullTemplate := fasttemplate.New(themeFullTemplateStr, "{{", "}}")
	instantiatedTheme := themeFullTemplate.ExecuteString(map[string]interface{}{"theme_optional_part": string(themeOptionalPart),
		"theme_mandatory_part": instantiatedThemeMandatoryPart})

	if _, err = tmpThemeFile.WriteString(instantiatedTheme); err != nil {
		err = eris.Wrap(internal.ConfigurationError, "Unable to write the theme file")

		return
	}

	res = tmpThemeFile.Name()

	return res, err
}

//nolint:funlen //Function contains the template
func AskUserToSelectCell(rows, cols uint32) (selectedRow, selectedCol uint32, err error) {
	var config internal.Configuration
	if err = viper.Unmarshal(&config); err != nil {
		internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

		return
	}

	position := config.GUIConfig.CellSelection.Position

	anchor, location, err := positionToAnchorAndLocation(position)
	if err != nil {
		err = eris.Wrapf(internal.ConfigurationError, "Invalid value of position '%s'", position)

		return
	}

	themeMandatoryPart :=
		`window {
	anchor:   {{window-anchor}};
	location: {{window-location}};
}

mainbox {
	children: [listview];
}

listview {
	columns: {{cols}};
	lines:   {{rows}};
}`

	themeMandatoryPartElements := map[string]interface{}{
		"cols":            fmt.Sprintf("%d", cols),
		"rows":            fmt.Sprintf("%d", rows),
		"window-anchor":   anchor,
		"window-location": location,
	}
	menuItems := linkedhashmap.New()
	for x := uint32(0); x < cols; x++ {
		for y := uint32(0); y < rows; y++ {
			menuItems.Put(strconv.Itoa(int((y*cols)+1+x)), int((y*cols)+1+x))
		}
	}

	_, selectedItemPayload, err := askUser(config.GUIConfig.CellSelection.ThemeFilePath,
		themeMandatoryPart, themeMandatoryPartElements, menuItems, "Select cell")
	if err != nil {
		err = eris.Wrapf(err, "Unable to ask for a cell to the user")

		return
	}

	selectedCell := selectedItemPayload.(int)

	selectedCol = (uint32(selectedCell) - 1) % cols
	selectedRow = (uint32(selectedCell) - 1) / cols

	return selectedRow, selectedCol, err
}

//nolint:funlen //Function contains the template
func AskUserToSelectGroup(predicate func(*windowmgt.WindowGroup) bool) (res *windowmgt.WindowGroup, err error) {
	var config internal.Configuration

	if err = viper.Unmarshal(&config); err != nil {
		internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

		return
	}

	position := config.GUIConfig.GroupSelection.Position

	anchor, location, err := positionToAnchorAndLocation(position)
	if err != nil {
		err = eris.Wrapf(internal.ConfigurationError, "Invalid value of position '%s'", position)

		return
	}

	themeMandatoryPart :=
		`window {
	anchor:   {{window-anchor}};
	location: {{window-location}};
}

mainbox {
  children: [inputbar,listview];
}

inputbar {
    children: [ prompt,entry ];
}`

	themeMandatoryPartElements := map[string]interface{}{
		"window-anchor":   anchor,
		"window-location": location,
	}

	wgm := windowmgt.GetWindowGroupManagerInstance()

	menuItems := linkedhashmap.New()

	for winGroupIt := wgm.GetGroups(predicate).Iterator(); winGroupIt.Next(); {
		curGroup := winGroupIt.Value().(*windowmgt.WindowGroup)
		menuItems.Put(curGroup.Name(), curGroup)
	}

	_, selectedItemPayload, err := askUser(config.GUIConfig.GroupSelection.ThemeFilePath,
		themeMandatoryPart, themeMandatoryPartElements, menuItems, "Select a group")
	if err != nil {
		err = eris.Wrapf(err, "Unable to ask for a cell to the user")

		return
	}

	res = selectedItemPayload.(*windowmgt.WindowGroup)

	return res, err
}

//nolint:funlen //Function contains the template
func AskUserToEnterNewGroup(prompt, placeholder string, options []string) (res string, err error) {
	var config internal.Configuration

	if err = viper.Unmarshal(&config); err != nil {
		internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

		return
	}

	position := config.GUIConfig.GroupCreation.Position

	anchor, location, err := positionToAnchorAndLocation(position)
	if err != nil {
		err = eris.Wrapf(internal.ConfigurationError, "Invalid value of position '%s'", position)

		return
	}

	themeMandatoryPart :=
		`mainbox {
    children:         [inputbar];
  }
  
  inputbar {
     children: [ prompt,entry ];
  }`

	themeMandatoryPartElements := map[string]interface{}{
		"window-anchor":   anchor,
		"window-location": location,
	}

	tmpThemeFile, err := generateThemeFile(config.GUIConfig.GroupSelection.ThemeFilePath, themeMandatoryPart, themeMandatoryPartElements)
	if err != nil {
		err = eris.Wrapf(err, "Unable to generate theme file for cell selection")

		return
	}

	defer os.Remove(tmpThemeFile)

	cmdToExec := fmt.Sprintf("rofi -i -dmenu -lines 0 -p %s %v -theme %s", shellescape.Quote(shellescape.StripUnsafe(prompt+":")),
		shellescape.QuoteCommand(options), tmpThemeFile)

	internal.NormalLogger.Debugf("Command to exec is '%s'", cmdToExec)
	c := clicmd.NewCommand(cmdToExec, clicmd.WithStandardStreams, clicmd.WithInheritedEnvironment(clicmd.EnvVars{}))

	if cmdErr := c.Execute(); cmdErr != nil {
		err = eris.Wrap(internal.InternalError, "Unable to execute rofi")

		return
	}

	if c.ExitCode() != 0 {
		err = eris.Wrapf(internal.InternalError, "Unable to execute the rofi command: '%s'", c.Stderr())

		return
	}

	res = strings.Trim(c.Stdout(), "\n")

	return res, err
}

//nolint:funlen //Function contains the template
func AskUserToSelectWindow(predicate func(*windowmgt.Window) bool) (res *windowmgt.Window, err error) {
	var config internal.Configuration

	if err = viper.Unmarshal(&config); err != nil {
		internal.NormalLogger.Errorf("Unable to load the configuration: %s", err)
		err = eris.Wrap(internal.InternalError, "Unable to load the configuration")

		return
	}

	position := config.GUIConfig.WindowSelection.Position

	anchor, location, err := positionToAnchorAndLocation(position)
	if err != nil {
		err = eris.Wrapf(internal.ConfigurationError, "Invalid value of position '%s'", position)

		return
	}

	themeMandatoryPart := `
	window {
		anchor:   {{window-anchor}};
		location: {{window-location}};
	}

	mainbox {
	  children: [inputbar,listview];
	}

	inputbar {
	    children: [ prompt,entry ];
	}`

	themeMandatoryPartElements := map[string]interface{}{
		"window-anchor":   anchor,
		"window-location": location,
	}

	wm, err := windowmgt.GetWindowManagerInstance()

	menuItems := linkedhashmap.New()

	for winIt := wm.GetWindows(predicate).Iterator(); winIt.Next(); {
		curWin := winIt.Value().(*windowmgt.Window)
		if predicate(curWin) {
			menuItems.Put(curWin.Name(), curWin)
		}
	}

	_, selectedItemPayload, err := askUser(config.GUIConfig.WindowSelection.ThemeFilePath,
		themeMandatoryPart, themeMandatoryPartElements, menuItems, "Select a window")
	if err != nil {
		err = eris.Wrapf(err, "Unable to ask for a cell to the user")
	}

	res = selectedItemPayload.(*windowmgt.Window)

	return res, err
}
