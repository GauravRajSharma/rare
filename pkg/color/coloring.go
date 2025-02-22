package color

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	escapeCode     = "\x1b"
	foregroundCode = "[3"
)

type ColorCode string

const (
	Reset         ColorCode = escapeCode + "[0m"
	Black                   = escapeCode + "[30m"
	Red                     = escapeCode + "[31m"
	Green                   = escapeCode + "[32m"
	Yellow                  = escapeCode + "[33m"
	Blue                    = escapeCode + "[34m"
	Magenta                 = escapeCode + "[35m"
	Cyan                    = escapeCode + "[36m"
	White                   = escapeCode + "[37m"
	BrightBlack             = escapeCode + "[30;1m"
	BrightRed               = escapeCode + "[31;1m"
	BrightGreen             = escapeCode + "[32;1m"
	BrightYellow            = escapeCode + "[33;1m"
	BrightBlue              = escapeCode + "[34;1m"
	BrightMagenta           = escapeCode + "[35;1m"
	BrightCyan              = escapeCode + "[36;1m"
	BrightWhite             = escapeCode + "[37;1m"

	Bold      = escapeCode + "[1m"
	Underline = escapeCode + "[4m"
)

// Enabled controls whether or not coloring is applied
var Enabled = true

var groupColors = [...]ColorCode{Red, Green, Yellow, Blue, Magenta, Cyan}

func init() {
	if runtime.GOOS == "windows" {
		Enabled = false
	} else if fi, err := os.Stdout.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			Enabled = false
		}
	}
}

// Wrap surroungs a string with a color (if enabled)
func Wrap(color ColorCode, s string) string {
	if !Enabled {
		return s
	}

	var sb strings.Builder
	sb.Grow(len(s) + 8)
	sb.WriteString(string(color))
	sb.WriteString(s)
	if len(s) < len(Reset) || s[len(s)-len(Reset):] != string(Reset) {
		sb.WriteString(string(Reset))
	}
	return sb.String()
}

func Wrapf(color ColorCode, s string, args ...interface{}) string {
	return Wrap(color, fmt.Sprintf(s, args...))
}

func Wrapi(color ColorCode, s interface{}) string {
	return Wrap(color, fmt.Sprintf("%v", s))
}

// WrapIndices color-codes by group pairs (regex-style)
//  [aStart, aEnd, bStart, bEnd...]
func WrapIndices(s string, groups []int) string {
	if !Enabled {
		return s
	}
	if len(groups) == 0 || len(groups)%2 != 0 {
		return s
	}

	var sb strings.Builder
	lastIndex := 0

	for i := 0; i < len(groups); i += 2 {
		start := groups[i]
		end := groups[i+1]
		color := groupColors[(i/2)%len(groupColors)]

		sb.WriteString(s[lastIndex:start])
		sb.WriteString(string(color))
		sb.WriteString(s[start:end])
		sb.WriteString(string(Reset))

		lastIndex = end
	}

	if lastIndex < len(s) {
		sb.WriteString(s[lastIndex:])
	}

	return sb.String()
}
