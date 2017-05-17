package cli

import (
	"os"

	"pkg.re/essentialkaos/ek.v8/arg"
	"pkg.re/essentialkaos/ek.v8/fmtc"
	"pkg.re/essentialkaos/ek.v8/knf"
	"pkg.re/essentialkaos/ek.v8/terminal"
	"pkg.re/essentialkaos/ek.v8/usage"

	"github.com/gongled/vgrepo/repo"
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
	KNF_STORAGE_URL  = "storage:url"
	KNF_STORAGE_PATH = "storage:path"
)

const (
	ARG_NO_COLOR = "nc:no-color"
	ARG_HELP     = "h:help"
	ARG_VER      = "v:version"
)

const (
	ERROR_UNSUPPORTED      = 1
	ERROR_INVALID_SETTINGS = 2
)

const CONFIG_FILE = "vgrepo.knf"

// ////////////////////////////////////////////////////////////////////////////////// //

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
	err := knf.Global(CONFIG_FILE)

	if err != nil {
		terminal.PrintErrorMessage(err.Error())
		os.Exit(ERROR_INVALID_SETTINGS)
	}
}

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
		terminal.PrintErrorMessage("Unknown command")
		os.Exit(ERROR_UNSUPPORTED)
	}

	if err != nil {
		terminal.PrintErrorMessage(err.Error())
		os.Exit(ERROR_UNSUPPORTED)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func addCommand(args []string) error {
	if len(args) != 3 {
		return fmtc.Errorf("Unable to handle %v arguments", len(args))
	}

	r := repo.NewRepository(
		knf.GetS(KNF_STORAGE_PATH),
		knf.GetS(KNF_STORAGE_URL),
		"openbox",
	)

	r.AddBox(args[0])

	fmtc.Println(r.Name)

	return nil
}

func deleteCommand(args []string) error {
	if len(args) != 1 {
		return fmtc.Errorf("Unable to handle %v arguments", len(args))
	} else {
		name := args[0]
		fmtc.Println(name)
	}

	return nil
}

func listCommand() error {
	return nil
}

func infoCommand(args []string) error {
	if len(args) != 1 {
		return fmtc.Errorf("Unable to handle %v arguments", len(args))
	} else {
		name := args[0]
		fmtc.Println(name)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func setUsageCommands(info *usage.Info) {
	info.AddCommand(CMD_ADD, "Add image to the Vagrant repository", "source", "name", "version")
	info.AddCommand(CMD_LIST, "Show the list of available images")
	info.AddCommand(CMD_DELETE, "Delete the image from the repository", "name", "version")
	info.AddCommand(CMD_INFO, "Display info of the particular repository", "name")
	info.AddCommand(CMD_HELP, "Display the current help message")
}

func setUsageOptions(info *usage.Info) {
	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")
}

func setUsageExamples(info *usage.Info) {
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
}

func showUsage() {
	info := usage.NewInfo(APP)

	setUsageCommands(info)
	setUsageOptions(info)
	setUsageExamples(info)

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
