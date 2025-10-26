package console

import (
	"fmt"
	"strings"
)

var lastLength int = 0

func setLastLength(text string) {
	lastNL := strings.LastIndex(text, "\n")
	lastLength = len(text[max(0, lastNL):])
}

func getSurrounding(text string) (string, string) {
	bs := strings.Repeat("\x08", lastLength)

	wslen := max(0, lastLength-len(text))
	ws := strings.Repeat(" ", wslen) + strings.Repeat("\x08", wslen)

	return bs, ws
}

func Printf(text string, placeholders ...any) {
	bs, ws := getSurrounding(text)
	total := fmt.Sprintf(text, placeholders...)

	fmt.Print(bs + total + ws)
	setLastLength(total)
}

func Print(text string) {
	bs, ws := getSurrounding(text)
	fmt.Print(bs + text + ws)

	setLastLength(text)
}

func Fprintln(text string, placeholders ...any) {
	_, ws := getSurrounding(text)
	total := text + ws + "\n"

	fmt.Printf(total, placeholders...)
	lastLength = 0
}

func Println(text string) {
	_, ws := getSurrounding(text)
	fmt.Println(text + ws)
	lastLength = 0
}
