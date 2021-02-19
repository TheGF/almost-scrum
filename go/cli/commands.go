package cli

import (
	"almost-scrum/core"
	"almost-scrum/web"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var config *core.Config

func usage() {
	fmt.Printf("usage: ash [-p <project-path>] [-u <user>] [-v] [-a] <command> [<args>]\n\n" +
		"These are the common Ash commands used in various situations.\n" +
		"\tinit              Initialize a project in the project path\n" +
		"\ttop [n]           Show top stories in current store\n" +
		"\tnew [title]       Create a task\n" +
		"\tedit [name]       Edit a task\n" +
		"\tdel [name]        Delete a task\n" +
		"\ttouch [name]      Focus on a task\n" +
		"\tmove [name]       Rename or move a task to a different board\n" +
		"\towner [name]      Assign the story to another user\n" +
		"\tcommit            Commit changes to the git repository\n" +
		"\tboard             List the boards and set the default\n" +
		"\tboard new <name>  Create a board with the provided name\n" +
		"\tusers add <id>    Add a user to current project\n" +
		"\tusers del <id>     Remove a user to current project\n" +
		"\tweb               Start the Web UI\n\n" +
		"\tserver <repo>     Start the Web UI as portal for projects in repo folder\n\n" +
		"\treindex [full]    Rebuild the search index \n\n" +
		"Options\n"+
		"\t-p <project-path> path where the current project is\n"+
		"\t-u <user>         impersonate a specific user (only for console client)\n"+
		"\t-v                shows verbose log\n"+
		"\t-x                exit when the user closes the browser (only for web UI)\n\n",
	)
}

var shortcuts = map[byte]string{
	'i': "init",
	't': "top",
	'n': "new",
	'e': "edit",
	'd': "del",
	'o': "owner",
	'c': "commit",
	'b': "board",
	'u': "users",
	'a': "a",
	'r': "reindex",
	's': "server",
}


func setLogLevel(logLevel string) {
	logLevel = strings.ToUpper(logLevel)
	switch logLevel {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	}
}

func replaceShortcuts(commands []string) []string {

	switch len(commands[0]) {
	case 1:
		if a := shortcuts[commands[0][0]]; a != "" {
			return append([]string{a}, commands[1:]...)
		}
	case 2:
		a := shortcuts[commands[0][0]]
		b := shortcuts[commands[0][1]]

		if a != "" && b != "" {
			return append([]string{a, b}, commands[1:]...)
		}
	}
	return commands

}

func matchGlobal(commands []string) bool {
	cmd := commands[0]
	if strings.HasSuffix(cmd, "!") {
		commands[0] = cmd[0:len(cmd)-1]
		return true
	} else {
		return false
	}

}

// ProcessArgs analyze the
func ProcessArgs() {
	var projectPath string
	var logLevel string
	var port string
	var username string
	var verbose bool
	var autoExit bool

	config = core.LoadConfig()

	flag.Usage = usage
	flag.StringVar(&projectPath, "p", ".",
		"Path where the project is. Default is current folder")

	flag.StringVar(&logLevel, "d", "error",
		"Log level to display")

	flag.StringVar(&username, "u", "",
		"Impersonate a different user")

	flag.StringVar(&port, "port", "8375",
		"HTTP port for the embedded web server")

	flag.BoolVar(&verbose, "v", false,
		"shows verbose log")

	flag.BoolVar(&autoExit, "x", false,
		"exit when the user closes the browser (only for web UI)")

	flag.Parse()

	if verbose {
		logLevel = "DEBUG"
	}
	setLogLevel(logLevel)

	if username != "" {
		core.SetSystemUser(username)
	}

	nArg := flag.NArg()
	if nArg == 0 {
		flag.Usage()
		return
	}

	commands := os.Args[len(os.Args)-nArg:]
	global := matchGlobal(commands)
	commands = replaceShortcuts(commands)
	switch commands[0] {
	case "init":
		processInit(projectPath, commands[1:])
	case "board":
		processBoard(projectPath, commands[1:])
	case "users":
		processUsers(projectPath, commands[1:])
	case "pwd":
		processPwd(commands[1:])
	case "top":
		processTop(projectPath, global, commands[1:])
	case "new":
		processNew(projectPath, commands[1:])
	case "edit":
		processEdit(projectPath, global, commands[1:])
	case "touch":
		processTouch(projectPath, global, commands[1:])
	case "owner":
		processOwner(projectPath, global, commands[1:])
	case "move":
		processMove(projectPath, global, commands[1:])
	case "commit":
		processCommit(projectPath, global)
	case "reindex":
		processReIndex(projectPath, commands[1:])
	case "web":
		web.StartWeb(projectPath, port, logLevel, autoExit, commands[1:])
	case "server":
		web.StartServer(port, logLevel, autoExit, commands[1:])

	default:
		logrus.Debugf("Unknown command %s", commands[0])
		flag.Usage()
	}

}
