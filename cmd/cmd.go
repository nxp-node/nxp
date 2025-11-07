package cmd

import (
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	"github.com/nxp-node/nxp/cmd/commands"
	"github.com/nxp-node/nxp/cmd/console"
	"golang.org/x/sys/windows"
)

var subcommands = map[string]Subcommand{
	"install": {
		Usage:         "[package]",
		ArgumentCount: MakeRange(1, 1),
		Description:   "installs the specified package",
		Function:      commands.Install,
	},
	"view": {
		Usage:         "[package]",
		ArgumentCount: MakeRange(1, 1),
		Description:   "views the specified package",
		Function:      commands.View,
	},
	"search": {
		Usage:         "[term]",
		ArgumentCount: MakeRange(1, -1),
		Description:   "searches the registry for all packages with the given term",
		Function: func(args []string) {
			println("not made yet")
		},
	},
}

func vtProcessing() {
	h, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		panic(err)
	}

	var mode uint32
	if err = windows.GetConsoleMode(h, &mode); err != nil {
		panic(err)
	}

	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if err = windows.SetConsoleMode(h, mode); err != nil {
		panic(err)
	}
}

func Main() {
	if runtime.GOOS == "windows" {
		vtProcessing()
	}

	if len(os.Args) <= 1 {
		viewCommands()
	} else {
		cmdName := os.Args[1]
		cmd, ok := subcommands[cmdName]

		if ok {
			args := slices.Delete(os.Args, 0, 2)
			length := uint16(len(args))

			var sRange string
			var maxSet bool

			if cmd.ArgumentCount.Maximum == nil {
				sRange = fmt.Sprintf("at least %d", cmd.ArgumentCount.Minimum)
				maxSet = false
			} else {
				sRange = fmt.Sprintf("%d-%d", cmd.ArgumentCount.Minimum, *cmd.ArgumentCount.Maximum)
				maxSet = true
			}

			if length < cmd.ArgumentCount.Minimum {
				console.Fprintln("%stoo few arguments (expected %s, got %d)", commands.Prefix, sRange, length)

				viewCommands()
				return
			} else if maxSet && length > *cmd.ArgumentCount.Maximum {
				console.Fprintln("%stoo many arguments (expected %s, got %d)", commands.Prefix, sRange, length)

				viewCommands()
				return
			} else {
				cmd.Function(args)
			}
		} else {
			console.Fprintln("%sunknown subcommand '%s'", commands.Prefix, cmdName)
			viewCommands()
		}
	}
}

func viewCommands() {
	console.Fprintln("%s[b]available subcommands:[/b]", commands.Prefix)

	maxName := 0
	maxUsage := 0

	for name, cmd := range subcommands {
		maxName = max(maxName, len(name))
		maxUsage = max(maxUsage, len(cmd.Usage))
	}

	for name, cmd := range subcommands {
		name += strings.Repeat(" ", maxName-len(name))

		if cmd.Usage != "" {
			name += " "

			name += fmt.Sprintf("[u]%s[/u]", cmd.Usage)
			name += strings.Repeat(" ", maxUsage-len(cmd.Usage))
		} else {
			name += strings.Repeat(" ", maxUsage+1)
		}

		console.Fprintln("%s> %s — %s", commands.Prefix, name, cmd.Description)
	}
}
