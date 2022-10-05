package core

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Fatal(format string, a ...any) {
	lineColor := color.New()
	lineColor.Add(color.FgRed)
	line := lineColor.Sprintf(format, a...)
	fmt.Println(line)
	os.Exit(-1)
}

func Info(format string, a ...any) {
	lineColor := color.New()
	lineColor.Add(color.FgRed)
	line := lineColor.Sprintf(format, a...)
	fmt.Println(line)
}
