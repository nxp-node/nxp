package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
	"strings"
)

var PrefixPre = "[#1e5688]nxp[/#1e5688] [#336997]"
var PrefixSuf = "[/#336997] "

var Prefix = PrefixPre + "│" + PrefixSuf

func init() {
	update()
}

var NXPConfig string

func update() {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("LOCALAPPDATA") + "/Olafcio Solutions/nxp"
	} else {
		dir = os.Getenv("HOME") + "/.config/Olafcio Solutions/nxp"
	}

	NXPConfig = dir
	os.MkdirAll(dir, 0700)

	path := dir + "/Accent Color"
	data, err := os.ReadFile(path)

	if errors.Is(err, fs.ErrNotExist) {
		os.WriteFile(path, []byte("#1e5688;#336997"), 0700)
	} else {
		var split = strings.SplitN(string(data[:]), ";", 2)

		PrefixPre = fmt.Sprintf("[%s]nxp[/%s] [%s]", split[0], split[0], split[1])
		PrefixSuf = fmt.Sprintf("[/%s] ", split[1])

		Prefix = PrefixPre + "│" + PrefixSuf
	}
}
