package commands

import (
	"os"

	"github.com/Olafcio1/nxp/cmd/console"
	"golang.org/x/term"
)

func overwriteCheck(path string) bool {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var data [1]byte
	if _, err := os.Stat(path); err == nil {
		console.Printf("nxp | '%s' already exists; overwrite (y = yes, n = cancel)? ", path)

		for {
			_, err := os.Stdin.Read(data[:])

			if err != nil {
				console.Println("\nnxp | accepted EOF; aborting")
				os.Exit(1)

				return true
			}

			if rune(data[0]) == 'y' {
				console.Println("y\nnxp | accepted 'y'; overwriting")
				break
			} else if rune(data[0]) == 'n' {
				console.Println("n\nnxp | accepted 'n'; aborting")
				os.Exit(0)

				return true
			}
		}

		return true
	} else {
		return false
	}
}
