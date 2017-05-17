package cli

import (
	"os"

	"pkg.re/essentialkaos/ek.v8/arg"
	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/usage"
	"pkg.re/essentialkaos/ek.v8/terminal"

	"github.com/gongled/vgrepo/prefs"
	"github.com/gongled/vgrepo/repo"
	//"github.com/gongled/vgrepo/meta"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "vgrepo"
	VER  = "2.0.0"
	DESC = "Simple CLI tool for managing Vagrant repositories"
)

const (
	CMD_ADD    = "add"
	CMD_DELETE = "delete"
	CMD_LIST   = "list"
	CMD_INFO   = "info"
	CMD_HELP   = "help"

	CMD_ADD_SHORTCUT    = "a"
	CMD_DELETE_SHORTCUT = "d"
	CMD_LIST_SHORTCUT   = "l"
	CMD_INFO_SHORTCUT   = "i"
)

const (
	ARG_NO_COLOR = "nc:no-color"
	ARG_HELP     = "h:help"
	ARG_VER      = "v:version"
)

const (
	ERROR_UNSUPPORTED = 0
)

const CONFIG_FILE = "vgrepo.knf"

// ////////////////////////////////////////////////////////////////////////////////// //

var preferences *prefs.Preferences

var argMap = arg.Map{
	ARG_NO_COLOR: {Type: arg.BOOL},
	ARG_HELP:     {Type: arg.BOOL, Alias: "u:usage"},
	ARG_VER:      {Type: arg.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
	args, errs := arg.Parse(argMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if arg.GetB(ARG_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if arg.GetB(ARG_VER) {
		showAbout()
		return
	}

	if arg.GetB(ARG_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	switch len(args) {
	case 0:
		showUsage()
		return
	case 1:
		processCommand(args[0], nil)
	default:
		processCommand(args[0], args[1:])
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func prepare() {
	preferences = prefs.New(CONFIG_FILE)
	errs := preferences.Validate()

	if len(errs) > 0 {
		for _, err := range errs {
			terminal.PrintErrorMessage(err.Error())
		}
		os.Exit(1)
	}
}

// process start source processing
func processCommand(cmd string, args []string) {
	var err error

	prepare()

	switch cmd {
	case CMD_ADD, CMD_ADD_SHORTCUT:
		err = addCommand(args)
	case CMD_DELETE, CMD_DELETE_SHORTCUT:
		err = deleteCommand(args)
	case CMD_LIST, CMD_LIST_SHORTCUT:
		err = listCommand()
	case CMD_INFO, CMD_INFO_SHORTCUT:
		err = infoCommand(args)
	case CMD_HELP:
		showUsage()
	default:
		fmtc.Printf("{r}Unknown command{!}\n")
		os.Exit(ERROR_UNSUPPORTED)
	}

	if err != nil {
		fmtc.Printf(err.Error())
		os.Exit(ERROR_UNSUPPORTED)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func addCommand(args []string) error {
	var err error

	if len(args) != 3 {
		err = fmtc.Errorf("Unable to handle %v arguments", len(args))
		return err
	}

	var (
		//src = args[0]
		name = args[1]
		//version = args[2]
	)

	r := repo.New(preferences, name)

	r.AddBox("qweqwe")

	return err
}

func deleteCommand(args []string) error {
	var err error

	if len(args) != 1 {
		err = fmtc.Errorf("Unable to handle %v arguments\n", len(args))
	} else {
		name := args[0]

		fmtc.Println(name)
	}

	return err
}

func listCommand() error {
	return nil
}

func infoCommand(args []string) error {
	var err error

	if len(args) != 1 {
		err = fmtc.Errorf("Unable to handle %v arguments\n", len(args))
	} else {
		name := args[0]
		fmtc.Println(name)
	}

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo(APP)

	info.AddCommand(CMD_ADD, "Add image to the Vagrant repository", "source", "name", "version")
	info.AddCommand(CMD_LIST, "Show the list of available images")
	info.AddCommand(CMD_DELETE, "Delete the image from the repository", "name", "version")
	info.AddCommand(CMD_INFO, "Display info of the particular repository", "name")
	info.AddCommand(CMD_HELP, "Display the current help message")

	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.AddExample(
		"add $HOME/powerbox-1.0.0.box powerbox 1.1.0",
		"Add image to the Vagrant repository",
	)

	info.AddExample(
		"list",
		"Show the list of available images",
	)

	info.AddExample(
		"remove powerbox 1.1.0",
		"Remove the image from the repository",
	)

	info.AddExample(
		"info powerbox",
		"Remove the image from the repository",
	)

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2014,
		Owner:   "Gleb E Goncharov",
		License: "MIT License",
	}

	about.Render()
}

// ////////////////////////////////////////////////////////////////////////////////// //
