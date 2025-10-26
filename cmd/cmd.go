package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Olafcio1/nxp/cmd/commands"
)

var subcommands = map[string]Subcommand{
	"install": {
		Usage:         "[package]",
		ArgumentCount: MakeRange(1, 1),
		Description:   "installs the specified package",
		Function:      commands.Install,
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

func Main() {
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
				fmt.Printf("nxp | too few arguments (expected %s, got %d)\n", sRange, length)

				viewCommands()
				return
			} else if maxSet && length > *cmd.ArgumentCount.Maximum {
				fmt.Printf("nxp | too many arguments (expected %s, got %d)\n", sRange, length)

				viewCommands()
				return
			} else {
				cmd.Function(args)
			}
		} else {
			fmt.Printf("nxp | unknown subcommand '%s'\n", cmdName)
			viewCommands()
		}
	}
}

func viewCommands() {
	fmt.Println("nxp | available subcommands:")

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

			name += cmd.Usage
			name += strings.Repeat(" ", maxUsage-len(cmd.Usage))
		} else {
			name += strings.Repeat(" ", maxUsage+1)
		}

		fmt.Printf("> %s — %s\n", name, cmd.Description)
	}
}
