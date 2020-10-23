package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// init in init_windows.go ensures that PowerShell and CMD emulate colors and handle
// other escape codes properly.
func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
