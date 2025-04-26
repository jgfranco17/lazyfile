package outputs

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrintColoredMessage(textColor string, message string, args ...any) string {
	var selectedColor color.Attribute
	switch strings.ToLower(textColor) {
	case "green":
		selectedColor = color.FgGreen
	case "yellow":
		selectedColor = color.FgYellow
	case "red":
		selectedColor = color.FgRed
	case "blue":
		selectedColor = color.FgBlue
	case "cyan":
		selectedColor = color.FgCyan
	default:
		selectedColor = color.FgWhite
	}
	colorFunc := color.New(selectedColor).SprintFunc()
	return colorFunc(fmt.Sprintf(message, args...))
}
