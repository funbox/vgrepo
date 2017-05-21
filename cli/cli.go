package cli

import (
	"os"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/fmtutil/table"
	"pkg.re/essentialkaos/ek.v9/knf"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/terminal"
	"pkg.re/essentialkaos/ek.v9/usage"

	"github.com/gongled/vgrepo/repo"
	"github.com/gongled/vgrepo/storage"
	"github.com/gongled/vgrepo/prefs"
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

const CONFIG_FILE = "/etc/vgrepo.knf"

// ////////////////////////////////////////////////////////////////////////////////// //

var optionsMap = options.Map{
	ARG_NO_COLOR: {Type: options.BOOL},
	ARG_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	ARG_VER:      {Type: options.BOOL, Alias: "ver"},
}

var preferences *prefs.Preferences

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
	opts, errs := options.Parse(optionsMap)

	if len(errs) != 0 {
		fmtc.Println("Arguments parsing errors:")

		for _, err := range errs {
			fmtc.Printf("  %s\n", err.Error())
		}

		os.Exit(1)
	}

	if options.GetB(ARG_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(ARG_VER) {
		showAbout()
		return
	}

	if options.GetB(ARG_HELP) || len(opts) == 0 {
		showUsage()
		return
	}

	switch len(opts) {
	case 0:
		showUsage()
		return
	case 1:
		processCommand(opts[0], nil)
	default:
		processCommand(opts[0], opts[1:])
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func prepare() {
	err := knf.Global(CONFIG_FILE)

	preferences = prefs.NewPreferences(
		knf.GetS(KNF_STORAGE_PATH),
		knf.GetS(KNF_STORAGE_URL),
	)

	if err != nil {
		terminal.PrintErrorMessage(err.Error())
		os.Exit(ERROR_INVALID_SETTINGS)
	}
}

func processCommand(cmd string, args []string) {
	prepare()

	switch cmd {
	case CMD_ADD, CMD_ADD_SHORTCUT:
		addCommand(args)
	case CMD_DELETE, CMD_DELETE_SHORTCUT:
		deleteCommand(args)
	case CMD_LIST, CMD_LIST_SHORTCUT:
		listCommand()
	case CMD_INFO, CMD_INFO_SHORTCUT:
		infoCommand(args)
	case CMD_HELP:
		showUsage()
	default:
		terminal.PrintErrorMessage("Error: unknown command %s", cmd)
		os.Exit(ERROR_UNSUPPORTED)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func addCommand(args []string) {
	if len(args) < 4 {
		terminal.PrintErrorMessage(
			"Error: unable to handle %v arguments",
			len(args),
		)
		os.Exit(1)
	}

	var (
		src      = args[0]
		name     = args[1]
		version  = args[2]
		provider = args[3]
	)

	r := repo.NewRepository(preferences, name)

	err := r.AddPackage(src, repo.NewPackage(name, version, provider))

	if err != nil {
		terminal.PrintErrorMessage("Error: %s", err.Error())
		os.Exit(1)
	}
}

func deleteCommand(args []string) {
	if len(args) < 2 {
		terminal.PrintErrorMessage("Error: name must be set")
		os.Exit(1)
	}

	var (
		name     = args[0]
		version  = args[1]
		provider = "virtualbox" // TODO: add support of partial remove by provider
	)

	r := repo.NewRepository(preferences, name)

	err := r.RemovePackage(repo.NewPackage(name, version, provider))

	if err != nil {
		terminal.PrintErrorMessage("Error: %s", err.Error())
		os.Exit(1)
	}
}

func listCommand() {
	s := storage.NewStorage(preferences)

	listTableRender(s.Repositories())
}

func infoCommand(args []string) {
	if len(args) < 1 {
		terminal.PrintErrorMessage("Error: name must be set")
		os.Exit(1)
	}

	name := args[0]

	infoTableRender(repo.NewRepository(preferences, name))
}

// ////////////////////////////////////////////////////////////////////////////////// //

func infoTableRender(r *repo.VRepository) {
	t := table.NewTable()
	table.HeaderCapitalize = true

	t.SetHeaders("Name", "Version", "Provider", "Checksum")
	t.SetAlignments(table.ALIGN_LEFT, table.ALIGN_CENTER, table.ALIGN_LEFT, table.ALIGN_LEFT)

	for _, v := range r.Versions {
		for _, p := range v.Providers {
			t.Add(r.Name, v.Version, p.Name, p.Checksum)
		}
	}

	if t.HasData() {
		t.Render()
	} else {
		terminal.PrintWarnMessage("Repository does not exist")
	}
}

func listTableRender(repos repo.VRepositoryList) {
	t := table.NewTable()
	table.HeaderCapitalize = true

	t.SetHeaders("Name", "Latest version", "URL")
	t.SetAlignments(table.ALIGN_LEFT, table.ALIGN_LEFT, table.ALIGN_LEFT)

	for _, r := range repos {
		if r.HasMeta() {
			t.Add(r.Name, r.LatestVersion().Version, r.URLMeta())
		}
	}

	if t.HasData() {
		t.Render()
	} else {
		terminal.PrintWarnMessage("No repositories yet")
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func setUsageCommands(info *usage.Info) {
	info.AddCommand(
		CMD_ADD,
		"Add image to the Vagrant repository",
		"source",
		"name",
		"version",
		"provider",
	)
	info.AddCommand(
		CMD_LIST,
		"Show the list of available images",
	)
	info.AddCommand(
		CMD_DELETE,
		"Delete the image from the repository",
		"name",
		"version",
		"provider",
	)
	info.AddCommand(
		CMD_INFO,
		"Display info of the particular repository",
		"name",
	)
	info.AddCommand(
		CMD_HELP,
		"Display the current help message",
	)
}

func setUsageOptions(info *usage.Info) {
	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")
}

func setUsageExamples(info *usage.Info) {
	info.AddExample(
		"add $HOME/powerbox-1.0.0.box powerbox 1.1.0 virtualbox",
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
