package console

import (
	"fmt"
	"strings"
)

func PrintEntries(target [][]string, linePrefix string) {
	table := entriesToTable(target)
	rowmaxs := calculateWidths(table)

	out := ""
	outlen := 0
	capped := 0

	for index, row := range table {
		for _, cell := range row {
			out += fmt.Sprintf(
				"%s[#676767]: [/#676767]%s",
				cell.key, cell.value,
			)

			outlen += int(cell.length)

			if outlen+4 > 84 {
				Println(linePrefix + out)

				out = ""
				capped = 0
			} else {
				out += strings.Repeat(" ", int(rowmaxs[index]))
			}

			capped += 1
		}
	}

	if out != "" {
		Println(linePrefix + out)
	}
}

type TableEntry struct {
	key    string
	value  string
	length uint
}

func entriesToTable(entries [][]string) [][]TableEntry {
	col := []TableEntry{}
	cols := [][]TableEntry{col}

	var total uint = 0
	var column uint = 0

	for _, tuple := range entries {
		key := tuple[0]
		value := tuple[1]

		length := uint(len(key) + len(value) + 2)
		newTotal := total + length

		entry := TableEntry{key, value, length}

		if newTotal > 84 {
			total = 0
			column += 1

			col = []TableEntry{entry}
			cols = append(cols, col)
		} else {
			col = append(col, entry)
		}
	}

	return cols
}

func calculateWidths(table [][]TableEntry) (widths []uint) {
	length := 0
	for _, rows := range table {
		length = max(length, len(rows))
	}

	for range length {
		widths = append(widths, 0)
	}

	for _, rows := range table {
		for index, row := range rows {
			widths[index] = max(
				widths[index],
				row.length,
			)
		}
	}

	return widths
}
