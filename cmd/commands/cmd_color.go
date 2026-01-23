package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nxp-node/nxp/cmd/console"
)

func Color(args []string) {
	if !strings.HasPrefix(args[0], "") || len(args[0]) != 7 {
		console.Printnln(Prefix + "error: the color argument must be a hex value, e.g. #ff0000")
		return
	}

	path := NXPConfig + "/Accent Color"

	primary := args[0]
	secondaryInt, err := strconv.ParseInt(args[0][1:], 16, 64)

	if err != nil {
		console.Printnln(Prefix + "error: failed to parse the hex value")
		console.Printnln(Prefix + "error: the color argument must be a hex value, e.g. #ff0000")
		return
	}

	secondaryR := (secondaryInt >> 16) & 255
	secondaryG := (secondaryInt >> 8) & 255
	secondaryB := secondaryInt & 255

	// Whoa this alghoritm works perfect!
	// Just removing 30 from the brightness, and, boom!
	// Equal colors!

	// (I used 50 because it's better to make a little darker)

	if secondaryR < secondaryG {
		if secondaryR < secondaryB {
			// smallest: R
			secondaryR -= 50
		} else {
			// smallest: B
			secondaryB -= 50
		}
	} else {
		if secondaryB < secondaryG {
			// smallest: B
			secondaryB -= 50
		} else {
			// smallest: G
			secondaryG -= 50
		}
	}

	secondary := "#" + fmt.Sprintf("%02x", max(0, secondaryR)) + fmt.Sprintf("%02x", max(0, secondaryG)) + fmt.Sprintf("%02x", max(0, secondaryB))
	os.WriteFile(path, []byte(primary+";"+secondary), 0700)

	update()

	console.Printnln(Prefix + "color updated")
}
