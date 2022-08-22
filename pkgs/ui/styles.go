package ui

import "github.com/charmbracelet/lipgloss"

const (
	ColorWhite = lipgloss.Color("#FFFFFF")
)

var (
	Bold = lipgloss.NewStyle().Bold(true).Foreground(ColorWhite).Render
)
