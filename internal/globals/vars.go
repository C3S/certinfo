package globals

import "github.com/fatih/color"

var (
	AllHosts   map[string]*Host
	Days       = 14
	Port       = 443
	IPversions = [2]int{4, 6}
	Timeout    = 5
	ShowErrors = true
	Blue       = color.New(color.FgBlue).SprintFunc()
	Magenta    = color.New(color.FgMagenta).SprintFunc()
	Green      = color.New(color.FgGreen).SprintFunc()
	Orange     = color.New(color.FgYellow).SprintFunc()
	Red        = color.New(color.FgRed).SprintFunc()
	// var Yellow = color.New(color.FgHiYellow).SprintFunc()
)
