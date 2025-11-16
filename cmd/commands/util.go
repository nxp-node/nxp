package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/nxp-node/nxp/cmd/console"
	"github.com/samber/lo"
	"golang.org/x/term"
)

const (
	OC_OVERWRITE int = iota
	OC_ABORTED
	OC_SAFE
)

func OC_OPT(index int) int {
	return index + 3
}

type Option struct {
	name        rune
	desc        string
	humanAction string
}

func overwriteCheck(path string, options ...Option) int {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	var data [1]byte
	if _, err := os.Stat(path); err == nil {
		console.Printf(
			Prefix+"'%s' already exists; overwrite (y = yes, n = cancel%s)? ",
			path,
			strings.Join(lo.Map(options, func(option Option, index int) string {
				return fmt.Sprintf(", %s = %s", string(option.name), option.desc)
			}), ""),
		)

		for {
			_, err := os.Stdin.Read(data[:])

			if err != nil {
				console.Fprintln("\n%saccepted EOF; aborting", Prefix)
				os.Exit(1)

				return OC_ABORTED
			}

			ch := rune(data[0])
			if ch == 'y' {
				console.Fprintln("y\n%saccepted 'y'; overwriting", Prefix)
				break
			} else if ch == 'n' {
				console.Fprintln("n\n%saccepted 'n'; aborting", Prefix)
				os.Exit(0)

				return OC_ABORTED
			} else {
				for i, option := range options {
					if option.name == ch {
						chStr := string(ch)
						console.Fprintln("%s\n%saccepted '%s'; %s", chStr, Prefix, chStr, option.humanAction)
						return OC_OPT(i)
					}
				}
			}
		}

		return OC_OVERWRITE
	} else {
		return OC_SAFE
	}
}
