package utils

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var ansi *regexp.Regexp = regexp.MustCompile("\x1b\\[(\\d+(?:;\\d+)*)m")

const (
	black int = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

const (
	foreground    = 30
	background    = 40
	bright        = 60
	bold          = 1
	faint         = 2
	italic        = 3
	underline     = 4
	strikethrough = 9
)

func colorName(code int) string {
	switch code {
	case black:
		return "black"
	case red:
		return "red"
	case green:
		return "green"
	case yellow:
		return "yellow"
	case blue:
		return "blue"
	case magenta:
		return "magenta"
	case cyan:
		return "cyan"
	case white:
		return "white"
	default:
		return ""
	}
}

func textModifier(code int) string {
	switch code {
	case bold:
		return "bold"
	case faint:
		return "faint"
	case italic:
		return "italic"
	case underline:
		return "underline"
	case strikethrough:
		return "strikethrough"
	default:
		return ""
	}
}

func attributesToClasses(attrList string) string {
	codes := strings.Split(attrList, ";")
	classNames := make([]string, 0, len(codes))

	for _, code := range codes {
		codeNum, err := strconv.Atoi(code)
		if err != nil {
			log.Printf("error parsing %s: %v", code, err)
			continue
		}

		className := ""
		if codeNum < 10 {
			className = textModifier(codeNum)
			if className != "" {
				classNames = append(classNames, className)
				continue
			}
		}

		// remove the last digit, so we get something like 30, 40, 90, or 100
		digit := codeNum % 10
		codeFamily := codeNum - digit

		// if the number is > 60, we're in a bright zone
		if codeFamily > bright {
			className = "bright-"
			codeFamily -= bright
		}

		if codeFamily == foreground {
			className += "fg-"
			color := colorName(digit)
			if color == "" {
				continue
			}
			className += color
			classNames = append(classNames, className)
			continue
		}

		if codeFamily == background {
			className += "bg-"
			color := colorName(digit)
			if color == "" {
				continue
			}
			className += color
			classNames = append(classNames, className)
			continue
		}
	}

	return strings.Join(classNames, " ")
}

func ANSItoHTML(s []byte) string {
	html := string(s)
	// convert new lines to line breaks
	html = strings.ReplaceAll(html, "\n", "<br/>\n")

	// convert reset codes to end spans
	html = strings.ReplaceAll(html, "\x1b[m", "</span>")
	html = strings.ReplaceAll(html, "\x1b[0m", "</span>")

	allMatches := ansi.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {
		html = strings.Replace(html, matches[0], fmt.Sprintf("<span class=%q>", attributesToClasses(matches[1])), 1)
	}

	return html
}
