package internal

import (
	"encoding/json"
	"os"
	"os/user"
	"path"

	"github.com/rotisserie/eris"
	"github.com/spf13/viper"
)

type FilterOperator string

type WindowProperty string

type GroupProperty string

const (
	WindowID    WindowProperty = "id"
	WindowClass                = "class"
	WindowName                 = "name"
	WindowType                 = "type"

	Or  FilterOperator = "or"
	And                = "and"

	GroupID   GroupProperty = "id"
	GroupName               = "name"
)

func (wf *WindowProperty) UnmarshalJSON(b []byte) (err error) {
	type WF WindowProperty
	var r *WF = (*WF)(wf)
	err = json.Unmarshal(b, &r)
	if err != nil {
		err = eris.Wrapf(InvalidArgumentError, "Invalid window field")
		return
	}
	switch *wf {
	case WindowID, WindowClass, WindowName, WindowType:
		return nil
	}
	err = eris.Wrapf(InvalidArgumentError, "Invalid window field")
	return
}

func (op *FilterOperator) UnmarshalJSON(b []byte) (err error) {
	type OP FilterOperator
	var r *OP = (*OP)(op)
	err = json.Unmarshal(b, &r)
	if err != nil {
		err = eris.Wrapf(InvalidArgumentError, "Invalid operator")
		return
	}
	switch *op {
	case And, Or:
		return nil
	}
	err = eris.Wrapf(InvalidArgumentError, "Invalid operator")
	return
}

// func NewWindowPropertyRegex(regexStr string) (res *WindowPropertyRegex, err error) {
// 	vanillaRegex, err := regexp.Compile(regexStr)
// 	if err != nil {
// 		err = eris.Wrapf(InvalidArgumentError, "Invalid window field")
// 		return
// 	}

// 	res = &WindowPropertyRegex{vanillaRegex}
// 	return
// }

const DefaultCellSelectionGUITheme string = `
  * {
    background-color: rgba(0,0,0,0);
  }

  window {
    fullscreen:       false;
    width:            150px;
    x-offset:         0px;
    y-offset:         0px;
    background-color: #282828;
    border-radius:    5px;
    padding:          5px 5px 5px 5px;
		font:             "Noto Sans Bold 12";
  }

  listview {
    padding: 0px 0px 0px 0px;
    spacing: 5px;
  }

  element {
    background-color: #404040;
    text-color:       #AAAAAA;
    border-radius:    5px;
    padding:          10px 0px 10px 0px;
    spacing:          5px;
  }

  element selected {
    background-color: #F0544C;
    text-color:       #EEEEEE;
  }

  element-text, element-icon {
    size: 0;
    horizontal-align: 0.5;
    vertical-align: 0.5;
  }
`
const DefaultGroupSelectionGUITheme string = `* {
  background-color: rgba(0,0,0,0);
}

window {
  fullscreen:       false;
  width:            400px;
  x-offset:         0px;
  y-offset:         0px;  
  background-color: #282828;
	padding:          0px 0px 0px 0px;
	font:             "Noto Sans 12";
}

entry {
	padding:           5px 5px 5px 5px;
  background-color:  #282828;
  text-color:        #AAAAAA;
  placeholder:       "Select a group";
  placeholder-color: #666666;
}

prompt {
  padding:          5px 5px 5px 5px;
  background-color: #F0544C;
	text-color:       #222222;
	font:             "Noto Sans Bold 12";
}

listview {
  fixed-height: false;
  dynamic:      true;
  scrollbar:    true;
	padding:      5px 5px 5px 5px;
	spacing:      5px;
}

scrollbar {   
  handle-color: #EEEEEE;
  handle-width: 5px;
}

element {
  background-color: #404040;
  text-color:       #AAAAAA;
  padding:          5px 5px 5px 5px;
}

element selected {
  background-color: #F0544C;
  text-color:       #EEEEEE;
}`

const DefaultWindowSelectionGUITheme string = `* {
  background-color: rgba(0,0,0,0);
}

window {
  fullscreen:       false;
  width:            400px;
  x-offset:         0px;
  y-offset:         0px;  
  background-color: #282828;
	padding:          0px 0px 0px 0px;
	font:             "Noto Sans 12";
}

entry {
	padding:           5px 5px 5px 5px;
  background-color:  #282828;
  text-color:        #AAAAAA;
  placeholder:       "Select a window";
  placeholder-color: #666666;
}

prompt {
  padding:          5px 5px 5px 5px;
  background-color: #F0544C;
	text-color:       #222222;
	font:             "Noto Sans Bold 12";
}

listview {
  fixed-height: false;
  dynamic:      true;
  scrollbar:    true;
	padding:      5px 5px 5px 5px;
	spacing:      5px;
}

scrollbar {   
  handle-color: #EEEEEE;
  handle-width: 5px;
}

element {
  background-color: #404040;
  text-color:       #AAAAAA;
  padding:          5px 5px 5px 5px;
}

element selected {
  background-color: #F0544C;
  text-color:       #EEEEEE;
}`

const DefaultGroupCreationGUITheme string = `* {
  background-color: rgba(0,0,0,0);
}

window {
  fullscreen:       false;
  width:            400px;
  x-offset:         0px;
  y-offset:         0px;  
  background-color: #282828;
	padding:          0px 0px 0px 0px;
	font:             "Noto Sans 12";
}

entry {
	padding:           5px 5px 5px 5px;
  background-color:  #282828;
  text-color:        #AAAAAA;
  placeholder:       "Enter the name of the new group";
  placeholder-color: #666666;
}

prompt {
  padding:          5px 5px 5px 5px;
  background-color: #F0544C;
	text-color:       #222222;
	font:             "Noto Sans Bold 12";
}`

type FilterConfiguration struct {
	WinProperty WindowProperty        `yaml:"window_property,omitempty" mapstructure:"window_property,omitempty"`
	Regex       string                `yaml:"regex,omitempty" mapstructure:"regex,omitempty"`
	Operator    FilterOperator        `yaml:"operator,omitempty" mapstructure:"operator,omitempty"`
	Filters     []FilterConfiguration `yaml:"filters,omitempty" mapstructure:"filters,omitempty"`
}

// type FilterConfiguration struct {
// 	WindowID    string `yaml:"windowid" mapstructure:"windowid"`
// 	Class       string `yaml:"class" mapstructure:"class"`
// 	Instance    string `yaml:"instance" mapstructure:"instance"`
// 	WindowRole  string `yaml:"window_role" mapstructure:"window_role"`
// 	WindowTitle string `yaml:"window_title" mapstructure:"window_title"`
// 	Workspace   string `yaml:"workspace" mapstructure:"workspace"`
// 	ConnMark    string `yaml:"conn_mark" mapstructure:"conn_mark"`
// }

type GroupConfiguration struct {
	ID      uint8                 `yaml:"id" mapstructure:"id"`
	Name    string                `yaml:"name" mapstructure:"name"`
	Color   string                `yaml:"color" mapstructure:"color"`
	Filters []FilterConfiguration `yaml:"filters" mapstructure:"filters"`
}

type GridConfiguration struct {
	Cols uint32 `yaml:"cols" mapstructure:"cols"`
	Rows uint32 `yaml:"rows" mapstructure:"rows"`
}

type GroupLabelsConfiguration struct {
	Font           string `yaml:"font" mapstructure:"font"`
	Position       string `yaml:"position" mapstructure:"position"`
	FontSize       uint32 `yaml:"font_size" mapstructure:"font_size"`
	DisplayTimeout uint32 `yaml:"display_timeout" mapstructure:"display_timeout"`
}

type SpecificGUIConfiguration struct {
	ThemeFilePath string `yaml:"theme_filepath" mapstructure:"theme_filepath"`
	Position      string `yaml:"position" mapstructure:"position"`
}

type GUIConfiguration struct {
	CellSelection   SpecificGUIConfiguration `yaml:"cellselection" mapstructure:"cellselection"`
	GroupSelection  SpecificGUIConfiguration `yaml:"groupselection" mapstructure:"groupselection"`
	WindowSelection SpecificGUIConfiguration `yaml:"windowselection" mapstructure:"windowselection"`
	GroupCreation   SpecificGUIConfiguration `yaml:"groupcreation" mapstructure:"groupcreation"`
}

type Configuration struct {
	GridConfig        GridConfiguration        `yaml:"grid" mapstructure:"grid"`
	GUIConfig         GUIConfiguration         `yaml:"gui" mapstructure:"gui"`
	GroupLabelsConfig GroupLabelsConfiguration `yaml:"group_labels" mapstructure:"group_labels"`
	GroupConfigs      []GroupConfiguration     `yaml:"groups" mapstructure:"groups"`
}

func InitConfigFile() (err error) {
	usr, err := user.Current()
	if err != nil {
		err = eris.Wrap(ConfigurationError, "Unable to get home directory of the user")

		return
	}

	configDirPath := path.Join(usr.HomeDir, ".config", "snapi3")
	configFilePath := path.Join(configDirPath, "snapi3.yml")

	viper.SetConfigName("snapi3")
	viper.SetConfigType("yml")
	viper.AddConfigPath(configDirPath)

	defaultConfig := Configuration{}
	defaultConfig.GridConfig.Cols = 3
	defaultConfig.GridConfig.Rows = 3
	defaultConfig.GroupLabelsConfig.Font = "Sans"
	defaultConfig.GroupLabelsConfig.FontSize = 10
	defaultConfig.GroupLabelsConfig.Position = "top"
	defaultConfig.GroupLabelsConfig.DisplayTimeout = 2

	defaultConfig.GUIConfig.CellSelection.Position = "center"
	defaultConfig.GUIConfig.CellSelection.ThemeFilePath = path.Join(configDirPath, "cellselection_gui_theme.rasi")

	defaultConfig.GUIConfig.GroupSelection.Position = "center"
	defaultConfig.GUIConfig.GroupSelection.ThemeFilePath = path.Join(configDirPath, "groupselection_gui_theme.rasi")

	defaultConfig.GUIConfig.WindowSelection.Position = "center"
	defaultConfig.GUIConfig.WindowSelection.ThemeFilePath = path.Join(configDirPath, "windowselection_gui_theme.rasi")

	defaultConfig.GUIConfig.GroupCreation.Position = "center"
	defaultConfig.GUIConfig.GroupCreation.ThemeFilePath = path.Join(configDirPath, "groupcreation_gui_theme.rasi")

	// groupConfig := GroupConfiguration{}
	// groupConfig.ID = 50
	// groupConfig.Color = "#AAAAAA"
	// groupConfig.Name = "toto"
	// filterConfig := FilterConfiguration{}
	// filterConfig.Operator = And

	// filterConfig.Regex = ".*"
	// filterConfig.WinProperty = WindowID
	// groupConfig.Filters = append(groupConfig.Filters, filterConfig)
	// defaultConfig.GroupConfigs = append(defaultConfig.GroupConfigs, groupConfig)

	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Colors.Background = "#404040"
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Colors.Text = "#AAAAAA"
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Colors.Placeholder = "#666666"
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.EntryConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Colors.Background = "#404040"
	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Colors.Text = "#AAAAAA"
	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.PromptConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Colors.Background = "#282828"
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Position = "center"
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.BorderRadius = 5
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.XOffset = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.YOffset = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Width = 500
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.NewGroupConfig.WindowConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Colors.Background = "#404040"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Colors.Text = "#AAAAAA"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Colors.Placeholder = "#666666"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.EntryConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Colors.Background = "#F0544C"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Colors.Text = "#AAAAAA"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.PromptConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Colors.SelectedBackground = "#F0544C"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Colors.SelectedText = "#EEEEEE"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Colors.Background = "#404040"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Colors.Text = "#AAAAAA"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.BorderRadius = 5
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Spacing = 5
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ItemConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.Colors.Handle = "#EEEEEE"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.HandleWidth = 5
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.ListViewConfig.Padding.Right = 0

	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Colors.Background = "#282828"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Position = "center"
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.BorderRadius = 5
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.XOffset = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.YOffset = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Width = 500
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Padding.Top = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Padding.Bottom = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Padding.Left = 0
	// defaultConfig.InterfaceConfig.SelectGroupConfig.WindowConfig.Padding.Right = 0

	viper.SetDefault("gui", defaultConfig.GUIConfig)
	viper.SetDefault("grid", defaultConfig.GridConfig)
	viper.SetDefault("group_labels", defaultConfig.GroupLabelsConfig)
	viper.SetDefault("groups", defaultConfig.GroupConfigs)

	if _, err = os.Stat(configDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, os.ModePerm)
		if err != nil {
			err = eris.Wrapf(InternalError, "Unable to create configuration directory '%s'", configDirPath)

			return
		}
	}

	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		if err = viper.SafeWriteConfig(); err != nil {
			err = eris.Wrap(ConfigurationError, "Unable to create the initial configuration")

			return
		}
	}

	if err = viper.ReadInConfig(); err != nil {
		VerboseLogger.Errorf(eris.ToString(err, true))
		err = eris.Wrap(ConfigurationError, "Unknown error")

		return
	}

	InitThemes()

	return err
}

func InitThemes() (err error) {

	usr, err := user.Current()
	if err != nil {
		err = eris.Wrap(ConfigurationError, "Unable to get home directory of the user")

		return
	}

	configDirPath := path.Join(usr.HomeDir, ".config", "snapi3")

	defaultCellSelectionGUIThemeFilePath := path.Join(configDirPath, "cellselection_gui_theme.rasi")
	if _, err = os.Stat(defaultCellSelectionGUIThemeFilePath); os.IsNotExist(err) {
		var defaultCellSelectionGUIThemeFile *os.File
		defaultCellSelectionGUIThemeFile, err = os.Create(defaultCellSelectionGUIThemeFilePath)
		_, err = defaultCellSelectionGUIThemeFile.WriteString(DefaultCellSelectionGUITheme)
		if err != nil {
			//TODO: Deal with the error
			return
		}
	}

	defaultGroupSelectionGUIThemeFilePath := path.Join(configDirPath, "groupselection_gui_theme.rasi")
	if _, err = os.Stat(defaultGroupSelectionGUIThemeFilePath); os.IsNotExist(err) {
		var defaultGroupSelectionGUIThemeFile *os.File
		defaultGroupSelectionGUIThemeFile, err = os.Create(defaultGroupSelectionGUIThemeFilePath)
		_, err = defaultGroupSelectionGUIThemeFile.WriteString(DefaultGroupSelectionGUITheme)
		if err != nil {
			//TODO: Deal with the error
			return
		}
	}

	defaultWindowSelectionGUIThemeFilePath := path.Join(configDirPath, "windowselection_gui_theme.rasi")
	if _, err = os.Stat(defaultWindowSelectionGUIThemeFilePath); os.IsNotExist(err) {
		var defaultWindowSelectionGUIThemeFile *os.File
		defaultWindowSelectionGUIThemeFile, err = os.Create(defaultWindowSelectionGUIThemeFilePath)
		_, err = defaultWindowSelectionGUIThemeFile.WriteString(DefaultWindowSelectionGUITheme)
		if err != nil {
			//TODO: Deal with the error
			return
		}
	}

	defaultGroupCreationGUIThemeFilePath := path.Join(configDirPath, "groupcreation_gui_theme.rasi")
	if _, err = os.Stat(defaultGroupCreationGUIThemeFilePath); os.IsNotExist(err) {
		var defaultGroupCreationGUIThemeFile *os.File
		defaultGroupCreationGUIThemeFile, err = os.Create(defaultGroupCreationGUIThemeFilePath)
		_, err = defaultGroupCreationGUIThemeFile.WriteString(DefaultGroupCreationGUITheme)
		if err != nil {
			//TODO: Deal with the error
			return
		}
	}

	return
}
