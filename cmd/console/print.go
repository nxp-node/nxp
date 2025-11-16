package console

import (
	"fmt"
	"strings"

	"github.com/iskaa02/qalam/bbcode"
)

var lastLength int = 0

func getLength(text string) int {
	length := 0
	for _, ch := range text {
		if ch == '\x1b' {
			for _, ch := range text {
				if ch == 'm' {
					break
				}
			}
		} else {
			length += 1
		}
	}

	return length
}

func setLastLength(text string) {
	lastNL := strings.LastIndex(text, "\n")
	lastLength = getLength(text[max(0, lastNL):])
}

func getSurrounding(text string) (string, string) {
	bs := strings.Repeat("\x08", lastLength)

	wslen := max(0, lastLength-getLength(text))
	ws := strings.Repeat(" ", wslen) + strings.Repeat("\x08", wslen)

	return bs, ws
}

func Printf(text string, placeholders ...any) {
	bs, ws := getSurrounding(text)
	bbcode.Printf(bs+text+ws, placeholders...)

	setLastLength(text)
}

func Print(text string) {
	bs, ws := getSurrounding(text)
	bbcode.Printf(bs + text + ws)

	setLastLength(text)
}

func Fprintln(text string, placeholders ...any) {
	_, ws := getSurrounding(text)
	total := text + ws + "\n"

	bbcode.Printf(total, placeholders...)
	lastLength = 0
}

func Println(text string) {
	_, ws := getSurrounding(text)

	display := bbcode.Sprintf(text)
	fmt.Println(display + ws)

	lastLength = 0
}

func Printnln(text string) {
	if lastLength != 0 {
		text = "\n" + text
	}

	Println(text)
}
