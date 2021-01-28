package main

import (
	"os"
	"path"

	"github.com/aacebedo/snapi3/internal"
	"github.com/aacebedo/snapi3/internal/hmi/cli"
	"github.com/aacebedo/snapi3/internal/hmi/gui"
	"github.com/aacebedo/snapi3/internal/windowmgt"
	"github.com/gofrs/flock"
	"github.com/rotisserie/eris"
	"github.com/spf13/cobra"
)

func CreateRootCmd() (res *cobra.Command) {
	//nolint:exhaustivestruct //cobra command struct has a lot of field that must stay defaulted
	res = &cobra.Command{
		Use:   "snapi3",
		Short: "A tool to quickly manipulate i3 floating windows",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}

	cli.AppendSnapCmd(res)
	cli.AppendCenterCmd(res)
	// cli.AppendGroupMgtCmd(res)
	// cli.AppendVisibilityMgtCmd(res)
	gui.AppendGUICmd(res)

	return
}

func main() {
	var weh *windowmgt.WindowEventHandler

	var wm *windowmgt.WindowManager

	var wgm *windowmgt.WindowGroupManager

	var ld *windowmgt.LabelDrawer

	var rootCmd *cobra.Command

	var windowGroupsVarFilePath string

	var err error

	fileLock := flock.New(path.Join(os.TempDir(), "snapi3.lock"))

	locked, err := fileLock.TryLock()
	if err != nil {
		goto exit
	}

	if !locked {
		err = eris.Wrap(internal.InternalError, "A instance of snapi3 is already running")

		goto exit
	}

	if err = internal.InitConfigFile(); err != nil {
		goto exit
	}

	internal.InitLoggers()

	ld, _ = windowmgt.GetLabelDrawerInstance()
	ld.Start()

	weh, _ = windowmgt.GetWindowEventHandlerInstance()
	weh.StartEventProcessing()

	if wm, err = windowmgt.GetWindowManagerInstance(); err != nil {
		goto exit
	}

	wgm = windowmgt.GetWindowGroupManagerInstance()

	if err = wm.LoadWindows(); err != nil {
		goto exit
	}

	if err = wgm.LoadWindowGroups(); err != nil {
		internal.NormalLogger.Warnf("Unable to load windows group from '%s', it will be recreated from scratch", windowGroupsVarFilePath)
	}

	rootCmd = CreateRootCmd()

	if err = rootCmd.Execute(); err != nil {
		goto exit
	}

	ld.Stop()
	err = wgm.SaveWindowGroups()
exit:
	if err == nil {
		os.Exit(0)
	} else {
		internal.NormalLogger.Error(eris.ToString(err, true))
		os.Exit(-1)
	}
}
