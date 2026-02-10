package ui

import "github.com/charmbracelet/lipgloss"

// Go quirk: can't use := at package level, need o do this
var (
	AppStyle = lipgloss.NewStyle().Margin(1, 2)

	// 2. The Key Data (Green for tz, Pink for user)
	KeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	ValStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF007C")).Bold(true)

	// 3. The Subtle Text (Gray for labels/help)
	SubtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

	// 4. Error (Red)
	ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
)
